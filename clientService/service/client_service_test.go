package service

import (
	"errors"
	"testing"
	"clientService/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockRequester struct {
	mock.Mock
}

func (m *MockRequester) Request(req data.RequestData) ([]data.Event, error) {
	args := m.Called(req)
	return args.Get(0).([]data.Event), args.Error(1)
}

func TestGetCriticalEvents_Success(t *testing.T) {
	mockRequester := new(MockRequester)

	events := []data.Event{
		{
			Criticality: 1,
			Timestamp:   "2025-01-01T00:00:00Z",
			EventMessage: "Critical event",
		},
	}

	requestData := data.RequestData{
		Limit:            5,
		CriticalityLevel: 5,
	}

	mockRequester.On("Request", requestData).Return(events, nil)

	es := NewEventService(mockRequester)

	result, err := es.GetCriticalEvents(5, 5)
	assert.NoError(t, err)
	assert.Equal(t, events, result)

	mockRequester.AssertCalled(t, "Request", requestData)
}

func TestGetCriticalEvents_RequestError(t *testing.T) {
	mockRequester := new(MockRequester)

	requestData := data.RequestData{
		Limit:            5,
		CriticalityLevel: 5,
	}

	mockRequester.On("Request", requestData).Return([]data.Event{}, errors.New("request error"))

	es := NewEventService(mockRequester)

	result, err := es.GetCriticalEvents(5, 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "request error")

	mockRequester.AssertCalled(t, "Request", requestData)
}

func TestGetCriticalEvents_InvalidLimit(t *testing.T) {
	mockRequester := new(MockRequester)

	es := NewEventService(mockRequester)

	result, err := es.GetCriticalEvents(0, 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "limit must be greater than 0")

	mockRequester.AssertNotCalled(t, "Request", mock.Anything)
}
