package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// NewWithdraw accepts request, save and return 200 if balance is enough.
func NewWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(userid).(int)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be 'application/json'", http.StatusBadRequest)
		return
	}

	newWD := storage.Withdraw{}
	if err := json.NewDecoder(r.Body).Decode(&newWD); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newWD.UserID = uid
	err := newWD.PostgresNewWD(ctx)

	switch err {
	case nil:
		w.WriteHeader(200) // approved
	case storage.ErrNoMoneyForWithdraw:
		w.WriteHeader(402) // deny
	default:
		http.Error(w, "can't accept new order", http.StatusInternalServerError)
	}
}

// GetWithdrawals returns list of withdrawals.
func GetWithdrawals(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(userid).(int)

	list, err := storage.PostgresGetWithdrawals(ctx, uid)
	if err != nil {
		http.Error(w, "can't get from DB", http.StatusInternalServerError)
	}
	if len(list) == 0 {
		w.WriteHeader(204)
		return
	}

	result, err := json.Marshal(list)
	if err != nil {
		log.Println("Error marshaling:", err)
		http.Error(w, "internal data error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}
