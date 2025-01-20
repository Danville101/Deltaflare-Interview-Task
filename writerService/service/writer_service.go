package service

import (
	"encoding/json"
	"errors"
	"log"
	"writerService/data"

	"github.com/nats-io/nats.go"
)

type RequestData struct {
	Limit            int    `json:"limit"`
	CriticalityLevel string `json:"criticalityLevel"`
}

type WriteService struct {
	natsConn       data.Subscriber
	influxDBClient data.Writer
}

func NewWriteService(nats data.Subscriber, influxDBClient data.Writer) (*WriteService, error) {
	return &WriteService{
		natsConn:       nats,
		influxDBClient: influxDBClient,
	}, nil
}

func (ws *WriteService) SubscribeAndProcess(topic string) error {
	_, err := ws.natsConn.Subscribe(topic, func(msg *nats.Msg) {
		var event data.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			return
		}

		if err := validateEvent(event); err != nil {
			log.Printf("Invalid event data: %v", err)
			return
		}


		err := ws.influxDBClient.WriteEvent(event)
		if err != nil {
			log.Printf("Error writing event to InfluxDB: %v", err)
		}
	})
	return err
}

func validateEvent(event data.Event) error {

	if event.Criticality < 0 || event.Criticality > 10 {
		return errors.New("criticality must be between 0 and 10")
	}

	
	if event.Timestamp == "" {
		return errors.New("timestamp cannot be empty")
	}



	if event.EventMessage == "" {
		return errors.New("eventMessage cannot be empty")
	}


	return nil
}
