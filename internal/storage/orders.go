package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	pgx "github.com/jackc/pgx/v5"
)

// PostgresNewOrder adds new order
func (a *Order) PostgresNewOrder(ctx context.Context) error {
	// write to postgres
	newOrder := `
	INSERT INTO orders(userid, number, status, accrual)
	VALUES ($1, $2, $3, 0)
;`
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
	}
	defer tx.Rollback(ctx)

	// search old number
	olduserid := 0
	err = db.QueryRow(ctx, "SELECT userid FROM orders WHERE userid =$1 AND number = $2", a.UserID, a.Number).Scan(&olduserid)
	switch err {
	case nil:
		// if same user
		if olduserid == a.UserID {
			return ErrConflictOrderUser
		}
		// if others
		return ErrConflictOrder
	case pgx.ErrNoRows:
	default:
		log.Println("postgres GET err: ", err)
		return err
	}

	// post
	_, er := tx.Exec(ctx, newOrder, a.UserID, a.Number, "NEW")
	if er != nil {
		errortext := er.Error()
		if errortext[len(errortext)-6:len(errortext)-1] == "23505" {
			log.Println("order exists: error 23505")
			return ErrConflictOrder
		}
		log.Println("2: post new order error:", err)
		return er
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
	}
	return nil
}

// PostgresGetOrder returns list of orders with all status by user
func PostgresGetOrder(ctx context.Context, uid int) ([]OrderToList, error) {
	var ordersList = `
	SELECT number, status, accrual, uploaded_at FROM orders 
	WHERE userid = $1
	ORDER BY uploaded_at ASC
	;`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
		return nil, fmt.Errorf("error with database")
	}
	defer tx.Rollback(ctx)

	rows, err := db.Query(ctx, ordersList, uid)
	if err != nil {
		fmt.Println("POSTGRES query error: ", err)
		return nil, fmt.Errorf("error with query database")
	}

	answer := make([]OrderToList, 0)

	tm := time.Now()

	for rows.Next() {
		New := OrderToList{}
		err := rows.Scan(&New.Number, &New.Status, &New.Accrual, &tm)
		if err != nil {
			fmt.Println("POSTGRES rows.Scan error: ", err)
		}
		New.Accrual = New.Accrual / 100
		New.UploadedAt = tm.Format(time.RFC3339)
		answer = append(answer, New)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
		return nil, fmt.Errorf("error with database")
	}
	return answer, nil
}
