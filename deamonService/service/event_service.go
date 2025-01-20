package services

import (
	"deamonService/data"
	"deamonService/models"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)





type EventService struct {
	publisher data.Publisher
}

func NewEventService(publisher data.Publisher) *EventService {
	return &EventService{publisher: publisher}
}

func (s *EventService) GenerateAndPublishEvent() {
	rand.Seed(time.Now().UnixNano())
	event := models.Event{
		Criticality:  rand.Intn(10) ,
		Timestamp:    time.Now().Format(time.RFC3339),
		EventMessage: "Random security event",
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event: %v", err)
		return
	}

	if err := s.publisher.Publish("events", data); err != nil {
		log.Printf("Error publishing event: %v", err)
	}
}
