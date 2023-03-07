// Package handler is operating endpoints
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(userid).(int)

	list, err := storage.PostgresGetOrder(ctx, uid)
	if err != nil {
		http.Error(w, "can't get from DB", http.StatusInternalServerError)
	}

	// if no results
	if len(list) == 0 {
		w.WriteHeader(204)
		return
	}

	// JSON
	result, err := json.Marshal(list)
	if err != nil {
		log.Println("Error marshaling:", err)
		http.Error(w, "internal data error", http.StatusInternalServerError)
	}
	log.Println("order:", string(result))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}
