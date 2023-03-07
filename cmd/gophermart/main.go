package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kormiltsev/gophermartbonus/internal/handler"
	"github.com/kormiltsev/gophermartbonus/internal/logger"
	"github.com/kormiltsev/gophermartbonus/internal/storage"
	"github.com/kormiltsev/gophermartbonus/internal/worker"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	blog := logger.NewLog()
	defer blog.Logger.Sync()

	undo := zap.RedirectStdLog(blog.Logger)
	defer undo()

	// blog := logger.NewLog()

	conf, err := storage.UploadConfigs()
	if err != nil {
		panic(err)
	}
	log.Println("service started on ", conf.Port)

	// connect PG
	conf.PostgresConnect(ctx)
	defer storage.PostgresClose()

	// upload data to RAM
	go storage.StartMemory(ctx)

	// start to work with external servise
	go worker.StartWorkers(ctx, conf)

	// chi router
	r := handler.NewRouter(conf)

	// set timeout
	server := &http.Server{
		Addr:         conf.Port,
		Handler:      http.TimeoutHandler(r, 50*time.Second, ""),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	fmt.Println(server.ListenAndServe())
}
