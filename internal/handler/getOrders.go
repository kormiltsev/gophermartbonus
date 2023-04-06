package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// GetOrders retruns all uploaded orders by user.
// @Tags 		Orders
// @Description Return actual user's list of orders with its statuses
// @Produce 	json
// @Success 	200 	{object} 	[]OrderToList
// @Success 	204 	{object}  	http.Response
// @Failure 	500 	{object}  	http.Response
// @Router 		/api/user/orders [get]
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
