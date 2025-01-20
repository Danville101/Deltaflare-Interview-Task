package service

import (
	"encoding/json"
	"errors"
	"testing"

	"readerService/data"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockReader struct {
	mock.Mock
}

func (m *MockReader) ReadDb(limit int, criticalityLevel int) ([]data.Event, error) {
	args := m.Called(limit, criticalityLevel)
	return args.Get(0).([]data.Event), args.Error(1)
}


type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(subject string, cb func(*nats.Msg)) (data.Subscription, error) {
	args := m.Called(subject, cb)

	if cb != nil {
	    req := RequestData{
		   Limit:            10,  
		   CriticalityLevel: 5, 
	    }
	    reqData, _ := json.Marshal(req)
 
	    msg := &nats.Msg{Data: reqData}
	    cb(msg) 
	}
 
	return args.Get(0).(data.Subscription), args.Error(1)
 }
 

func TestSubscribeAndProcess_Success(t *testing.T) {
	mockReader := new(MockReader)
	mockSubscriber := new(MockSubscriber)

	events := []data.Event{
		{Criticality: 1, Timestamp: "2025-01-01T00:00:00Z", EventMessage: "Test event"},
	}
	mockReader.On("ReadDb", 10, 5).Return(events, nil)

	mockSubscription := new(MockSubscription)
	mockSubscriber.On("Subscribe", "query.event", mock.Anything).Return(mockSubscription, nil)

	rs, err := NewReaderService(mockSubscriber, mockReader)
	assert.NoError(t, err)

	err = rs.SubscribeAndProcess()
	assert.NoError(t, err)
	mockReader.AssertCalled(t, "ReadDb", 10, 5)
}

func TestSubscribeAndProcess_ReadDbError(t *testing.T) {
	mockReader := new(MockReader)
	mockSubscriber := new(MockSubscriber)

	mockReader.On("ReadDb", 10, 5).Return([]data.Event{}, errors.New("database error"))

	mockSubscription := new(MockSubscription)
	mockSubscriber.On("Subscribe", "query.event", mock.Anything).Return(mockSubscription, nil)

	rs, err := NewReaderService(mockSubscriber, mockReader)
	assert.NoError(t, err)

	err = rs.SubscribeAndProcess()
	assert.NoError(t, err)
	mockReader.AssertCalled(t, "ReadDb", 10, 5)
}

type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) Unsubscribe() error {
	args := m.Called()
	return args.Error(0)
}
