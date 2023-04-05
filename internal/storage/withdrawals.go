package storage

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
)

// UpdateUsers update catalog in RAM. Run in case:
//
// start service (with old data) or
//
// new data from external service or
//
// new success withdraw
func (r *RAM) UpdateUsers(ctx context.Context) {

	var usersSum = `
	SELECT accrual.userid, COALESCE(accruals, 0), COALESCE(withdrawals,0) FROM
    ((SELECT userid, COALESCE(SUM(accrual),0) as accruals FROM orders
	GROUP BY userid) AS accrual
    LEFT JOIN (SELECT userid, COALESCE(SUM(sum),0) as withdrawals FROM withdrawals
	GROUP BY userid) AS w ON accrual.userid = w.userid );
	`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
	}
	defer tx.Rollback(ctx)

	rows, err := db.Query(ctx, usersSum)
	if err != nil {
		log.Println("POSTGRES query sum error: ", err)
	}

	cat.mu.Lock()
	for rows.Next() {

		new := User{}
		err := rows.Scan(&new.UserID, &new.Sum, &new.Withdrawsum)
		if err != nil {
			log.Println("POSTGRES rows.Scan error: ", err)
		}

		new.Withdrawsum = new.Withdrawsum / 100
		new.Sum = math.Round(new.Sum-new.Withdrawsum*100) / 100

		cat.catalog[new.UserID] = new
	}
	cat.mu.Unlock()
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
	}
	// log.Println("cat = ", cat)
}

// PostgresNewWD push new wd to postgres, returns "no money" error if not enough money.
func (wd *Withdraw) PostgresNewWD(ctx context.Context) error {
	us, ok := cat.catalog[wd.UserID]
	if us.Sum-us.Withdrawsum < wd.Sum && ok {
		return ErrNoMoneyForWithdraw
	}

	var inserWD = `
	INSERT INTO withdrawals(userid, ordernumber, sum, processed_at)
	SELECT $1, $2, $3, NOW()
	WHERE EXISTS (SELECT FROM orders
			WHERE (SELECT COALESCE(SUM(o.accrual),0)-COALESCE(SUM(w.sum),0)
		FROM orders AS o
		LEFT JOIN withdrawals AS w ON o.userid = w.userid
		WHERE o.userid = $4
		GROUP BY o.userid) > $5)
		ON CONFLICT(ordernumber) DO NOTHING
		;`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := db.Exec(ctx, inserWD, wd.UserID, wd.Order, wd.Sum*100, wd.UserID, wd.Sum*100)
	if err != nil {
		log.Println("POSTGRES query sum error: ", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
		return err
	}

	if rows.RowsAffected() == 0 {
		return ErrNoMoneyForWithdraw
	}
	// ctx Background cose run after request will be over?
	go cat.UpdateUsers(context.Background())
	return nil
}

// PostgresGetWithdrawals returns list of withdrawals by user.
func PostgresGetWithdrawals(ctx context.Context, uid int) ([]WithdrawList, error) {
	var wdList = `
	SELECT ordernumber, sum, processed_at FROM withdrawals 
	WHERE userid = $1
	ORDER BY processed_at ASC
	;`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
		return nil, fmt.Errorf("error with database")
	}
	defer tx.Rollback(ctx)

	rows, err := db.Query(ctx, wdList, uid)
	if err != nil {
		fmt.Println("POSTGRES query wd error: ", err)
		return nil, fmt.Errorf("error with wd query database")
	}

	answer := make([]WithdrawList, 0)
	tm := time.Now()
	for rows.Next() {
		New := WithdrawList{}
		err := rows.Scan(&New.Order, &New.Sum, &tm)
		if err != nil {
			fmt.Println("POSTGRES rows.Scan error: ", err)
		}
		New.Sum = New.Sum / 100
		New.ProcessedAt = tm.Format(time.RFC3339)
		answer = append(answer, New)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
		return nil, fmt.Errorf("error with database")
	}
	return answer, nil
}
