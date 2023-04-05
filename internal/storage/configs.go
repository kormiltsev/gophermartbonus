package storage

import (
	"errors"
	"flag"
	"log"
	"strings"

	env "github.com/caarlos0/env/v6"
)

// ServerConfigs is server settings.
type ServerConfigs struct {
	Port            string `env:"RUN_ADDRESS"`            // RUN_ADDRESS or -a
	UserEndpoint    string `env:"-"`                      // /api/user
	Register        string `env:"-"`                      // /register
	Login           string `env:"-"`                      // /login
	UserUpload      string `env:"-"`                      // /orders
	Balance         string `env:"-"`                      // /balance
	AskWithdraw     string `env:"-"`                      // /balance/withdraw
	Withdrawals     string `env:"-"`                      // /withdrawals
	ExternalService string `env:"ACCRUAL_SYSTEM_ADDRESS"` // /api/orders/{number} ACCRUAL_SYSTEM_ADDRESS or -r
	DBURI           string `env:"DATABASE_URI"`           // DATABASE_URI or -d
}

// UploadConfigs parse flags end ENV.
func UploadConfigs() (*ServerConfigs, error) {
	conf := ServerConfigs{
		UserEndpoint: "/api/user",
		Register:     "/register",
		Login:        "/login",
		UserUpload:   "/orders",
		Balance:      "/balance",
		AskWithdraw:  "/balance/withdraw",
		Withdrawals:  "/withdrawals",
	}
	// reading flags
	conf.Flags()

	// check environment RUN_ADDRESS, ACCRUAL_SYSTEM_ADDRESS, DATABASE_URI. Top priority.
	conf.Environment()

	// env is top priority, if not null set it up
	if conf.Port == "" {
		conf.Port = "localhost:8080"
	}

	// conf.ExternalService = "http://localhost:8081"
	if conf.ExternalService == "" {
		log.Println("no external service parameters, using default http://localhost:8081")
		conf.ExternalService = "http://localhost:8081"
		// return nil, errors.New("no external service parameters")
	}

	if conf.DBURI == "" {
		log.Println("no database parameters")
		return nil, errors.New("no database parameters")
	}

	// in case `env:"SERVER_ADDRESS"` = something long
	conf.port()
	return &conf, nil
}

// Environment returns ENV values.
func (c *ServerConfigs) Environment() {
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("got from env:", c)
}

// Flags returns service parameters in case of flags.
func (c *ServerConfigs) Flags() {
	// Server conf flags
	port := flag.String("a", "localhost:8080", "service port")
	extra := flag.String("r", "", "bonus calculation service adress")
	storage := flag.String("d", "", "DB connection link")

	flag.Parse()

	c.Port = *port
	c.ExternalService = *extra
	c.DBURI = *storage
}

// port check for correct port data.
func (c *ServerConfigs) port() {
	a := strings.Split(c.Port, ":")
	switch len(a) {
	case 0:
		c.Port = ":8080"
	case 1:
		c.Port = ":" + a[0]
	case 2:
		c.Port = a[0] + ":" + a[1]
	default:
		c.Port = (a[1])[2:] + ":" + a[2]
	}
}
