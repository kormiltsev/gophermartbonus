package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
	"github.com/theplant/luhn"
)

// NewOrder accepts new order by user and push to BD.
func NewOrder(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(userid).(int)

	simplebody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order := string(simplebody)
	defer r.Body.Close()

	if len(order) == 0 {
		http.Error(w, "empty body in post request", http.StatusBadRequest)
		return
	}

	// check for Luhn
	num, err := strconv.Atoi(order)
	if err != nil {
		log.Println("order number to int fail")
		http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
	}

	if !luhn.Valid(num) {
		// if !encode.LuhnValid(order) {
		log.Println("luhn fail")
		http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
		w.WriteHeader(422)
		return
	}

	neworder := storage.Order{
		UserID: uid,
		Number: order,
	}
	log.Println("new order: ", neworder)

	err = neworder.PostgresNewOrder(ctx)
	switch err {
	case nil:
		w.WriteHeader(202) // new order created
	case storage.ErrConflictOrder:
		w.WriteHeader(409) // order uploaded by other user
	case storage.ErrConflictOrderUser:
		w.WriteHeader(200) // order uploaded
	default:
		http.Error(w, "can't accept new order", http.StatusInternalServerError)
	}
}
