package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// Balance returns actual user's balance.
// @Tags 		Balance
// @Description Return actual user's balance
// @Accept  	any
// @Produce 	json
// @Success 	200 	{object} 	storage.User
// @Failure 	500 	{object}  	http.Response
// @Router 		/api/user/balance [get]
func Balance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(userid).(int)

	userstat := storage.User{
		UserID: uid,
	}

	// ask postgres
	err := userstat.PostgresGetBalance(ctx)
	if err != nil {
		http.Error(w, "can't get balance from DB", http.StatusInternalServerError)
	}

	// JSON
	result, err := json.Marshal(userstat)
	if err != nil {
		log.Println("Error marshaling:", err)
		http.Error(w, "internal data error", http.StatusInternalServerError)
	}

	log.Println("balanse:", string(result))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}
