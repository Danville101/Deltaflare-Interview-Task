package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"writerService/data"
	"writerService/service"
	"github.com/joho/godotenv"
	"fmt"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	    log.Fatalf("Error loading .env file")
	}
	host := os.Getenv("NATSHOST")
	port := os.Getenv("NATSPORT")
	nastUser := os.Getenv("NATSUSER")
	natsPassword := os.Getenv("NATSPASSWORD")


	url := fmt.Sprintf("nats://%s:%s@%s:%s",nastUser,natsPassword, host, port)


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


	eventService, err := service.NewWriteService(subscriber, influxDBClient)
	if err != nil {
		log.Fatalf("Error initializing EventService: %v", err)
	}


	if err := eventService.SubscribeAndProcess("events"); err != nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")
}
