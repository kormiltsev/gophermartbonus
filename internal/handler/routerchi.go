// Package handler is operating endpoints
package handler

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// NewRouter manages endpoints
func NewRouter(con *storage.ServerConfigs) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	// r.Use(middleware.Timeout(5000 * time.Millisecond))

	// for users
	r.Route(con.UserEndpoint, func(r chi.Router) { //con.UserEndpoint
		// login or register
		r.Post(con.Register, NewUser) //con.Register
		r.Post(con.Login, LoginUser)

		r.Route("/", func(r chi.Router) {
			// get user id from cookies
			r.Use(Authorization)

			r.Post(con.UserUpload, NewOrder)
			r.Get(con.UserUpload, GetOrders)
			r.Get(con.Balance, Balance)
			r.Post(con.AskWithdraw, NewWithdraw)
			r.Get(con.Withdrawals, GetWithdrawals)
		})
	})
	return r
}
