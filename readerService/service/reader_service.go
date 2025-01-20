package service

import (
	"encoding/json"
	"log"

	"readerService/data"
	"github.com/nats-io/nats.go"
)

type RequestData struct {
	Limit            int    `json:"limit"`
	CriticalityLevel int `json:"criticalityLevel"`
}

type ReaderService struct {
	natsConn       data.Subscriber
	influxDBClient data.Reader
}

func NewReaderService(nats data.Subscriber, influxDBClient data.Reader) (*ReaderService, error) {
	return &ReaderService{
		natsConn:       nats,
		influxDBClient: influxDBClient,
	}, nil
}

func (rs *ReaderService) SubscribeAndProcess() error {
	_, err := rs.natsConn.Subscribe("query.event", func(m *nats.Msg) {
		var req RequestData
		if err := json.Unmarshal(m.Data, &req); err != nil {
			log.Printf("Error unmarshaling request data: %v", err)
			m.Respond([]byte("error unmarshaling request data"))
			return
		}

		if req.Limit <= 0 {
			m.Respond([]byte("limit must be greater than zero"))
			return
		}

		events, err := rs.influxDBClient.ReadDb(req.Limit, req.CriticalityLevel)
		if err != nil {
			log.Printf("Failed to get critical events: %v", err)
			m.Respond([]byte("failed to get critical events"))
			return
		}

		responseData, err := json.Marshal(events)
		if err != nil {
			log.Printf("Error marshaling events: %v", err)
			m.Respond([]byte("error marshaling events"))
			return
		}

		if err := m.Respond(responseData); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	})

	if err != nil {
		return err
	}

	return nil
}
