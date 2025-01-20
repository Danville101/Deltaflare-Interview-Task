package data
import (
	"github.com/nats-io/nats.go"
)

type Publisher interface{
	Publish(subject string , data []byte) error
}



type EventPublisher struct{
	nc *nats.Conn

}


func NewEventPublisher(nc *nats.Conn) *EventPublisher{
	return &EventPublisher{nc: nc}
}


func(ep *EventPublisher) Publish(subject string , data []byte) error{
	return ep.nc.Publish(subject, data)
}