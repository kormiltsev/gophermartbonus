package storage

import "time"

// Order used by worker to update status.
type Order struct {
	UserID     int       `json:"-"`
	Number     string    `json:"order"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"-"`
}

// User used in authorization and balance counting.
type User struct {
	UserID      int     `json:"-"`
	Login       string  `json:"-"`
	Pass        string  `json:"-"`
	Sum         float64 `json:"current"`
	Withdrawsum float64 `json:"withdrawn"`
}

// Withdraw is row in withdrawals.
type Withdraw struct {
	UserID      int       `json:"-"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"-"`
}

// WithdrawList is list of withdrawals by user.
type WithdrawList struct {
	UserID      int     `json:"-"`
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

// OrderToList used to make list of orders uploaded.
type OrderToList struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
