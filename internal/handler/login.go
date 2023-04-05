package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kormiltsev/gophermartbonus/internal/encode"
	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// Userauth is login-password format from user.
type Userauth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginUser accepts login-password and add Bearer.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("loginuser")
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be 'application/json'", http.StatusBadRequest)
		return
	}

	newuser := Userauth{}
	if err := json.NewDecoder(r.Body).Decode(&newuser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hide password in hash
	k := sha256.Sum256([]byte(newuser.Login + newuser.Password))
	passhash := hex.EncodeToString(k[:])
	// =====================

	newusertostorage := storage.User{
		Login: newuser.Login,
		Pass:  passhash,
	}

	// check exists?
	uid, err := newusertostorage.PostgresLoginUser(ctx)
	if err != nil {
		log.Println("can't find user:", err)
		http.Error(w, "user not exists", http.StatusUnauthorized)
		return
	}

	// if error and user id 0
	if uid < 1 {
		log.Println("wrong iser id from DB")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// create header Auth
	text, err := encode.Shifu(uid)
	if err != nil {
		log.Println("can't encode:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Create a Bearer
	var bearer = "Bearer " + text
	w.Header().Add("Authorization", bearer)

	w.WriteHeader(200)

}

// NewUser is for register new user.
func NewUser(w http.ResponseWriter, r *http.Request) {
	log.Println("newuser")
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be 'application/json'", http.StatusBadRequest)
		return
	}

	newuser := Userauth{}
	if err := json.NewDecoder(r.Body).Decode(&newuser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hide password in hash
	k := sha256.Sum256([]byte(newuser.Login + newuser.Password))
	passhash := hex.EncodeToString(k[:])
	// =====================

	newusertostorage := storage.User{
		Login: newuser.Login,
		Pass:  passhash,
	}

	// add new user to PG
	uid, err := newusertostorage.PostgresNewUser(ctx)
	if err != nil {
		log.Println("can't add new user:", err)
		http.Error(w, "login exists", http.StatusConflict)
	}

	// create header Auth
	text, err := encode.Shifu(uid)
	if err != nil {
		log.Println("can't encode:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	// Create a Bearer
	var bearer = "Bearer " + text
	w.Header().Add("Authorization", bearer)

	w.WriteHeader(200)
}
