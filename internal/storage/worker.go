package storage

import (
	"context"
	"fmt"
	"log"
)

// UpdateRows post workers results and return new array same lenght
func UpdateRows(ctx context.Context, list []Order, QtyToSend int) ([]Order, error) {

	sqlStatement := `
		UPDATE orders
		SET status = $1, accrual = $2, uploaded_at = NOW()
		WHERE number = $3;`

	sqlStatementreq := `
		SELECT id, number, status FROM orders
		WHERE id > $1 
		AND (status ='PROCESSING' OR status = 'REGISTERED' OR status = 'NEW')
		ORDER BY id
		LIMIT $2;`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
		return nil, fmt.Errorf("can't touch DB")
	}
	defer tx.Rollback(ctx)

	if len(list) != 0 {
		for i := 0; i < len(list); i++ {
			_, err := db.Exec(ctx, sqlStatement, list[i].Status, list[i].Accrual*100, list[i].Number)
			if err != nil {
				log.Println("2.1: error push to postgres: ", err)
			}
		}
	}

	rows, err := db.Query(ctx, sqlStatementreq, currentID, QtyToSend)
	if err != nil {
		fmt.Println("POSTGRES query error: ", err)
		return nil, fmt.Errorf("can't upload from DB")
	}
	answer := make([]Order, 0)
	for rows.Next() {
		New := Order{}
		err := rows.Scan(&currentID, &New.Number, &New.Status)
		if err != nil {
			fmt.Println("POSTGRES rows.Scan error: ", err)
		}
		answer = append(answer, New)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
		return nil, fmt.Errorf("can't upload from DB")
	}
	go cat.UpdateUsers(ctx)
	return answer, nil
}
