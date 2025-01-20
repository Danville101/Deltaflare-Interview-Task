package service

import (
	"errors"
	"log"
	"clientService/data"

)

type EventService interface {
	GetCriticalEvents(limit int, criticalityLevel int) ([]data.Event, error)
}

type EventServiceImpl struct {
	requester data.Requester
}

func NewEventService(requester data.Requester) *EventServiceImpl {
	return &EventServiceImpl{requester: requester}
}

func (es *EventServiceImpl) GetCriticalEvents(limit int, criticalityLevel int) ([]data.Event, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}


	if criticalityLevel < 0 || criticalityLevel > 10 {
		return nil, errors.New("criticalityLevel must be between 0 and 10")
	}

	req := data.RequestData{
		Limit:            limit,
		CriticalityLevel: criticalityLevel,
	}

	events, err := es.requester.Request(req)
	if err != nil {
		log.Printf("Error fetching critical events: %v", err)
		return nil, err
	}


	return events, nil
}

