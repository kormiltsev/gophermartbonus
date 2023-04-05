package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/kormiltsev/gophermartbonus/internal/encode"
	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// usera using as context value in request.
type usera string

// userid using as context value in request.
var userid usera

// proxyHandle decode Bearer and return user id in context.
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uid := 0
		var err error

		// get Auth bearer
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)

		if authHeader == "" || reqToken == authHeader {
			log.Println("no header Authorization bearer")
			http.Error(w, "Authentication header not present or malformed", http.StatusUnauthorized)
			return
		}

		// decode user id
		uid, err = encode.Deshifu(reqToken)
		if err != nil {
			log.Println("wrong:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// if id exists in DB
		u := storage.User{UserID: uid}
		err = u.PostgresUserID(ctx)
		if err != nil {
			log.Println("user not found by id:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// code user id
		val, err := encode.Shifu(u.UserID)
		if err != nil {
			log.Println("can't encode:", err)
		}

		// Create a Bearer
		var bearer = "Bearer " + val
		w.Header().Add("Authorization", bearer)

		// context vith value
		ctx = context.WithValue(r.Context(), userid, uid) //u.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
