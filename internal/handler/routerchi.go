// Package handler is operating endpoints.
package handler

import (
	_ "net/http/pprof"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// NewRouter manages endpoints.
func NewRouter(con *storage.ServerConfigs) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	// r.Use(middleware.Timeout(5000 * time.Millisecond))

	r.Mount("/debug", Profiler())

	// for users
	r.Route(con.UserEndpoint, func(r chi.Router) {
		// login or register
		r.Post(con.Register, NewUser)
		r.Post(con.Login, LoginUser)

		r.Route("/", func(r chi.Router) {
			// get user id from Bearer
			r.Use(Authorization) // check bearer auth

			r.Post(con.UserUpload, NewOrder)       // post new order
			r.Get(con.UserUpload, GetOrders)       // get orders list
			r.Get(con.Balance, Balance)            // get actual balance
			r.Post(con.AskWithdraw, NewWithdraw)   // request withdrawal
			r.Get(con.Withdrawals, GetWithdrawals) // get list of withdrawals
		})
	})
	return r
}
