package service

import (
	"encoding/json"
	"errors"
	"testing"

	"writerService/data"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) WriteEvent(event data.Event) error {
	args := m.Called(event)
	return args.Error(0)
}


type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(subject string, cb func(*nats.Msg)) (data.Subscription, error) {
	args := m.Called(subject, cb)
	return args.Get(0).(data.Subscription), args.Error(1)
}

type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) Unsubscribe() error {
	args := m.Called()
	return args.Error(0)
}

func TestSubscribeAndProcess_Success(t *testing.T) {
	mockWriter := new(MockWriter)
	mockSubscriber := new(MockSubscriber)

	event := data.Event{
		Criticality: 1,
		Timestamp:   "2025-01-01T00:00:00Z",
		EventMessage: "Test event",
	}
	eventData, _ := json.Marshal(event)

	mockWriter.On("WriteEvent", event).Return(nil)

	mockSubscription := new(MockSubscription)
	mockSubscriber.On("Subscribe", "test.topic", mock.Anything).Return(mockSubscription, nil).Run(func(args mock.Arguments) {
		callback := args.Get(1).(func(*nats.Msg))
		callback(&nats.Msg{Data: eventData})
	})

	ws, err := NewWriteService(mockSubscriber, mockWriter)
	assert.NoError(t, err)

	err = ws.SubscribeAndProcess("test.topic")
	assert.NoError(t, err)

	mockWriter.AssertCalled(t, "WriteEvent", event)
}

func TestSubscribeAndProcess_WriteEventError(t *testing.T) {
	mockWriter := new(MockWriter)
	mockSubscriber := new(MockSubscriber)

	event := data.Event{
		Criticality: 1,
		Timestamp:   "2025-01-01T00:00:00Z",
		EventMessage: "Test event",
	}
	eventData, _ := json.Marshal(event)

	mockWriter.On("WriteEvent", event).Return(errors.New("write error"))

	mockSubscription := new(MockSubscription)
	mockSubscriber.On("Subscribe", "test.topic", mock.Anything).Return(mockSubscription, nil).Run(func(args mock.Arguments) {
		callback := args.Get(1).(func(*nats.Msg))
		callback(&nats.Msg{Data: eventData})
	})

	ws, err := NewWriteService(mockSubscriber, mockWriter)
	assert.NoError(t, err)

	err = ws.SubscribeAndProcess("test.topic")
	assert.NoError(t, err)

	mockWriter.AssertCalled(t, "WriteEvent", event)
}

func TestSubscribeAndProcess_UnmarshalError(t *testing.T) {
	mockWriter := new(MockWriter)
	mockSubscriber := new(MockSubscriber)

	invalidData := []byte("invalid data")

	mockSubscription := new(MockSubscription)
	mockSubscriber.On("Subscribe", "test.topic", mock.Anything).Return(mockSubscription, nil).Run(func(args mock.Arguments) {
		callback := args.Get(1).(func(*nats.Msg))
		callback(&nats.Msg{Data: invalidData})
	})

	ws, err := NewWriteService(mockSubscriber, mockWriter)
	assert.NoError(t, err)

	err = ws.SubscribeAndProcess("test.topic")
	assert.NoError(t, err)

	mockWriter.AssertNotCalled(t, "WriteEvent", mock.Anything)
}
