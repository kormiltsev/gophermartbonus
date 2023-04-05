package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	chi "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

var operates atomic.Bool
var answerType int

func main() {
	answerType = flags()
	// go timer()
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/api/orders/{number}", Reply)
		// 	r.Use(AuthMiddleware)
		r.Get("/h", hw)
	})
	log.Println("starts blackbox on 8081\ntype case: ", answerType)
	http.ListenAndServe("localhost:8081", r)
}

func timer() {
	for {
		time.Sleep(5 * time.Second)
		operates.Store(true)
		log.Println("ban activated")
		time.Sleep(2 * time.Second)
		operates.Store(false)
		log.Println("ban released")
	}
}

// answer
type tp struct {
	Order   string
	Status  string
	Accrual float64
}

func Reply(w http.ResponseWriter, r *http.Request) {
	// time delay
	time.Sleep(400 * time.Millisecond)

	if operates.Load() {
		log.Println("got request, but BAN is active")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Retry-After", "2")
		w.WriteHeader(429)
		w.Write([]byte("No more than N requests per minute allowed"))
		return
	}

	asknum := chi.URLParam(r, "number")
	log.Print("order number:", asknum)

	switch answerType {
	case 1:
		// PROCCESSED every time, sum is 123.89
		completeStatusForEverybody(asknum, w, r)
	default:

		// rundom status answer
		randomStatusByTime(asknum, w, r)
	}

}

func AuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func hw(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Pong"))
}

func completeStatusForEverybody(asknum string, w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(asknum)
	ans := tp{}
	if err != nil {
		ans = tp{
			Order:  asknum,
			Status: "INVALID",
		}
	} else {
		ans = tp{
			Order:   asknum,
			Status:  "PROCESSED",
			Accrual: 123.89,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	result, err := json.Marshal(ans)
	if err != nil {
		log.Println("error in marshal", err)
	}
	w.Write(result)
}

func randomStatusByTime(asknum string, w http.ResponseWriter, r *http.Request) {
	a := time.Now().UnixNano()

	_, err := strconv.Atoi(asknum)
	if err != nil {
		ans := tp{
			Order:  asknum,
			Status: "INVALID",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		result, err := json.Marshal(ans)
		if err != nil {
			log.Println("error in marshal", err)
		}
		w.Write(result)
		return
	}

	switch (a / 1000) % 10 {
	case 2:

		ans := tp{
			Order:  asknum,
			Status: "REGISTERED",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		result, err := json.Marshal(ans)
		if err != nil {
			log.Println("error in marshal", err)
		}
		w.Write(result)
		return
	case 3:

		ans := tp{
			Order:   asknum,
			Status:  "PROCESSED",
			Accrual: 500.32,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		result, err := json.Marshal(ans)
		if err != nil {
			log.Println("error in marshal", err)
		}
		w.Write(result)
		return
	case 4:

		ans := tp{
			Order:  asknum,
			Status: "INVALID",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		result, err := json.Marshal(ans)
		if err != nil {
			log.Println("error in marshal", err)
		}
		w.Write(result)
		return
	default:

		ans := tp{
			Order:  asknum,
			Status: "PROCESSING",
		}
		result, err := json.Marshal(ans)
		if err != nil {
			log.Println("cant marshal ", ans, err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(result)
		return
	}
}

// Flags returns parameters in case of flags
func flags() int {
	// Server conf flags
	port := flag.Int("a", 1, "type of answers")

	flag.Parse()

	return *port
}
