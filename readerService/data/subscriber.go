package data

import (
	"github.com/nats-io/nats.go"
)


type Message interface {
	GetData() []byte
	GetSubject() string
}

type Subscription interface {
	Unsubscribe() error
}

type Subscriber interface {
	Subscribe(subject string, cb func(*nats.Msg)) (Subscription, error)
}

type EventSubscriber struct {
	nc *nats.Conn
}


func NewEventSubscriber(nc *nats.Conn) *EventSubscriber {
	return &EventSubscriber{nc: nc}
}

type messageWrapper struct {
	msg *nats.Msg
}

func (m *messageWrapper) GetData() []byte {
	return m.msg.Data
}

func (m *messageWrapper) GetSubject() string {
	return m.msg.Subject
}


type subscriptionWrapper struct {
	sub *nats.Subscription
}

func (s *subscriptionWrapper) Unsubscribe() error {
	return s.sub.Unsubscribe()
}

func (ep *EventSubscriber) Subscribe(subject string, cb func(*nats.Msg)) (Subscription, error) {
	sub, err := ep.nc.Subscribe(subject, cb)
	if err != nil {
		return nil, err
	}
	return &subscriptionWrapper{sub: sub}, nil
}
