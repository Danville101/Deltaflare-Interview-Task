package data

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)


type RequestData struct {
	Limit            int `json:"limit"`
	CriticalityLevel  int `json:"criticalityLevel"`
}

type Event struct {
	Criticality int    `json:"criticality"`
	Timestamp   string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

type Requester interface {
	Request(req RequestData) ([]Event, error)
}

type ClientPublisher struct {
	nc *nats.Conn
}

func NewClientPublisher(nc *nats.Conn) *ClientPublisher {
	return &ClientPublisher{nc: nc}
}


func (cp *ClientPublisher) ClosePublisher(){
	cp.nc.Close()
}

func (cp *ClientPublisher) Request(req RequestData) ([]Event, error) {
	subject := "query.event"


	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}


	msg, err := cp.nc.Request(subject, data, 2*time.Second)
	if err != nil {
		return nil, err
	}


	var events []Event
	err = json.Unmarshal(msg.Data, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
