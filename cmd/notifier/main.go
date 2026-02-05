package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"lms_backend/internal/notifier/worker"
	"lms_backend/pkg/broker"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	eventBroker := broker.NewEventBroker(redisAddr)
	defer eventBroker.Close()

	notificationWorker := worker.NewNotificationWorker(eventBroker, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := notificationWorker.Run(ctx); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	log.Println("Notifier Service running...")
	<-sigChan
	log.Println("Shutting down...")
}
