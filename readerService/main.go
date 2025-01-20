package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"readerService/data"
	"readerService/service"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	natshost := os.Getenv("NATSHOST")
	natsport := os.Getenv("NATSPORT")
	url := fmt.Sprintf("nats://%s:%s", natshost, natsport)

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	subscriber := data.NewEventSubscriber(nc)

	influxURL := os.Getenv("INFLUXURL")
	org := os.Getenv("INFLUXDB_ORG")
	bucket := os.Getenv("INFLUXDB_BUCKET")
	token := os.Getenv("INFLUXTOKEN")

	influxDBClient := data.NewInfluxDBClient(influxURL, token, org, bucket)

	eventService, err := service.NewReaderService(subscriber, influxDBClient)
	if err != nil {
		log.Fatalf("Error initializing ReaderService: %v", err)
	}

	if err := eventService.SubscribeAndProcess(); err != nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")
}
