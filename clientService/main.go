package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"clientService/data"
	"clientService/service"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/nats-io/nats.go"
	"strconv"
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
 	criticalityLevel := os.Getenv("CRITICALITYLEVEL")
	 criticalityLevelInt, err := strconv.Atoi(criticalityLevel)
	if err !=nil{
		log.Println("Error converting criticalityLevelreq", err)
		
	}

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}


	publisher := data.NewClientPublisher(nc)
	eventService := service.NewEventService(publisher)
	defer  publisher.ClosePublisher()



	data, err:= eventService.GetCriticalEvents(10, criticalityLevelInt )
	if err !=nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}else{
		for i, v := range data{
			fmt.Printf(
				"Event %d :{\n" +
				"    Criticality: %d\n" +
				"    Timestamp: %s\n" +
				"    EventMessage: %s\n" +
				"}\n",
				i+1, v.Criticality, v.Timestamp, v.EventMessage,
			 )
			 
		}
		
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	

	log.Println("\n Shutting down gracefully...")
}
