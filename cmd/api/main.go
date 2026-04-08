package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soundmarket/backend/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("app init failed: %v", err)
	}
	defer application.DB.Close()
	defer application.Redis.Close()

	server := &http.Server{
		Addr:              ":" + application.Config.AppPort,
		Handler:           application.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("soundmarket api started on :%s", application.Config.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
}
