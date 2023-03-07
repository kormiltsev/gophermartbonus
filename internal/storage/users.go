package storage

import (
	"context"
	"errors"
	"log"
)

// PostgresLoginUser returns user id or error if not found or wrong password
func (u *User) PostgresLoginUser(ctx context.Context) (int, error) {
	uid := 0
	err := db.QueryRow(ctx, "SELECT id FROM users WHERE login = $1 AND pass = $2", u.Login, u.Pass).Scan(&uid)
	if err != nil {
		log.Println("PG: wrond login-pass", err)
		return 0, ErrConflictLoginUser
	}
	return uid, nil
}

// PostgresNewUser add new user
func (u *User) PostgresNewUser(ctx context.Context) (int, error) {
	// write to postgres
	newOrder := `
	INSERT INTO users(login, pass, sum, withdrawsum, created_at)
	VALUES ($1, $2, $3, $4, NOW())
	RETURNING id
;`
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("1: begin error:", err)
	}
	defer tx.Rollback(ctx)

	// post
	uid := 0
	er := db.QueryRow(ctx, newOrder, u.Login, u.Pass, 0, 0).Scan(&uid)
	if er != nil {
		errortext := er.Error()
		if errortext[len(errortext)-6:len(errortext)-1] == "23505" {
			log.Println("user exists: error 23505")
			return 0, ErrConflictNewUser
		}
		log.Println("post new user error:", err)
		return 0, er
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("3: commit err: ", err)
	}
	return uid, nil
}

// InternalUser check for exists
func (u *User) PostgresUserID(ctx context.Context) error {
	uid := 0
	err := db.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", u.UserID).Scan(&uid)
	if err != nil {
		log.Println("PG: user not found", err)
		return errors.New("user not found")
	}
	return nil
}
