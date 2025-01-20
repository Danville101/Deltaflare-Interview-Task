// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"deamonService/data"
	"deamonService/service"
	"github.com/joho/godotenv"
	"fmt"
)

func main() {

	err := godotenv.Load()
	if err != nil {
	    log.Fatalf("Error loading .env file")
	}
	host := os.Getenv("NATSHOST")
	port := os.Getenv("NATSPORT")
	url := fmt.Sprintf("nats://%s:%s", host, port)
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	publisher := data.NewEventPublisher(nc)
	eventService := services.NewEventService(publisher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
		cancel()
	}()

	run(ctx, eventService)
}

func run(ctx context.Context, eventService *services.EventService) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			return
		case <-ticker.C:
			eventService.GenerateAndPublishEvent()
		}
	}
}
