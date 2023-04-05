package storage

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RAM is used to get balance in memory to quick access. Updates every time of order status changed or successful withdrawal.
type RAM struct {
	mu      *sync.Mutex
	catalog map[int]User
}

// Creating type RAM.
var cat = RAM{
	mu:      &sync.Mutex{},
	catalog: make(map[int]User, 0),
}

// Postgress DB in here.
var db *pgxpool.Pool

// Error templates.
var (
	ErrConflictOrder      = errors.New(`conflict: order already exists`)
	ErrConflictOrderUser  = errors.New(`conflict: order already uploaded`)
	ErrConflictNewUser    = errors.New(`conflict: user already exists`)
	ErrNoMoneyForWithdraw = errors.New(`not enough money or order already exists`)
	ErrConflictLoginUser  = errors.New(`wrong login/password`)
)

// currentID used by workers in querying.
var currentID = 0

// PostgresConnect make connection with DB or panic.
func (c *ServerConfigs) PostgresConnect(ctx context.Context) {
	// connect to DB
	poolConfig, err := pgxpool.ParseConfig(c.DBURI)
	if err != nil {
		log.Panic("Unable to parse database_url:", err)
	}
	log.Println(poolConfig)

	db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Panic("Unable to create connection pool:", err)
	}
	// create tables if not exists. Orders table
	var orders = `
			CREATE TABLE IF NOT EXISTS orders(
				id serial primary key,
			  userid INTEGER not null,
			  number TEXT not null unique,
			  status TEXT not null,
			  accrual INTEGER,
			  uploaded_at TIMESTAMPTZ DEFAULT Now()
			);
		  `
	log.Println("creation table or")
	_, err = db.Exec(ctx, orders)
	if err != nil {
		log.Println("error in create tble orders, table exists?", err)
	}

	// users table
	var users = `
	CREATE TABLE IF NOT EXISTS users(
		id serial primary key,
	  login VARCHAR(128) not null unique,
	  pass VARCHAR(128) not null,
	  sum INTEGER,
	  withdrawsum INTEGER,
	  created_at TIMESTAMPTZ DEFAULT Now()
	);
  `
	_, err = db.Exec(ctx, users)
	log.Println("creation table usr")
	if err != nil {
		log.Println("error in create tble users, table exists?", err)
	}

	// withdraws table
	var withdraws = `
	CREATE TABLE IF NOT EXISTS withdrawals(
	id serial primary key,
	userid INTEGER,
	ordernumber VARCHAR(64) not null unique,
	sum INTEGER,
	processed_at TIMESTAMPTZ DEFAULT Now()
	);
	  `
	_, err = db.Exec(ctx, withdraws)
	log.Println("creation table wd")
	if err != nil {
		log.Println("error in create tble withdrawals, table exists?, ", err)
	}

	if !PostgresPing(ctx) {
		panic("postgres not connected")
	}
}

// PostgresPing returns TRUE if DB available.
func PostgresPing(ctx context.Context) bool {
	_, err := db.Exec(ctx, ";")
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// PostgresClose close all connections.
func PostgresClose() {
	db.Close()
}

// PostgresGetBalance returns sum of orders with status PROCCESSED and sum all withdrawals.
func (u *User) PostgresGetBalance(ctx context.Context) error {
	us, ok := cat.catalog[u.UserID]
	if !ok {
		u.Sum = 0
		u.Withdrawsum = 0
		return nil
	}
	u.Sum = us.Sum
	u.Withdrawsum = us.Withdrawsum
	return nil
}

// StartMemory to keep users info im memory.
func StartMemory(ctx context.Context) {
	cat.UpdateUsers(ctx)
}
