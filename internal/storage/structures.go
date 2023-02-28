package storage

import "time"

type Order struct {
	UserID     int       `json:"-"`
	Number     string    `json:"order"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"-"`
}

type User struct {
	UserID      int     `json:"-"`
	Login       string  `json:"-"`
	Pass        string  `json:"-"`
	Sum         float64 `json:"current"`
	Withdrawsum float64 `json:"withdrawn"`
}

type Withdraw struct {
	UserID      int       `json:"-"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"-"`
}

type WithdrawList struct {
	UserID      int     `json:"-"`
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type OrderToList struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
