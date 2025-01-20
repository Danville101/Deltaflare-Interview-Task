package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"deamonService/models"
)

type MockPublisher struct {
	mock.Mock
	PublishedData map[string][]byte
}

func (m *MockPublisher) Publish(topic string, data []byte) error {
	args := m.Called(topic, data)
	m.PublishedData = map[string][]byte{
		topic: data,
	}
	return args.Error(0)
}

func TestGenerateAndPublishEvent(t *testing.T) {
	mockPublisher := new(MockPublisher)
	eventService := NewEventService(mockPublisher)


	mockPublisher.On("Publish", "events", mock.Anything).Return(nil)


	eventService.GenerateAndPublishEvent()


	mockPublisher.AssertExpectations(t)


	publishedData := mockPublisher.PublishedData["events"]
	assert.NotNil(t, publishedData, "Expected data to be published to the 'events' topic")

	var event models.Event
	err := json.Unmarshal(publishedData, &event)
	assert.NoError(t, err, "Error unmarshalling event data")

	assert.Equal(t, "Random security event", event.EventMessage, "Expected EventMessage to be 'Random security event'")
	assert.True(t, event.Criticality >= 0 && event.Criticality <= 10, "Expected Criticality to be between 0 and 10")
	_, err = time.Parse(time.RFC3339, event.Timestamp)
	assert.NoError(t, err, "Expected Timestamp to be in RFC3339 format")
}
