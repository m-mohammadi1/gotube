package channelauthhandler

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gotube/internal/config"
	"gotube/pkg/model"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// MockChannelRepository is a mock implementation of ChannelRepository.
type MockChannelRepository struct {
	Channels map[int64]model.Channel
	mutex    sync.Mutex
	nextID   int64
}

// NewMockChannelRepository creates a new instance of MockChannelRepository.
func NewMockChannelRepository() *MockChannelRepository {
	return &MockChannelRepository{
		Channels: make(map[int64]model.Channel),
		nextID:   1,
	}
}

// UpdateOrCreate updates or creates a channel in the repository.
func (m *MockChannelRepository) UpdateOrCreate(ctx context.Context, channel model.Channel) (*model.Channel, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Assign a unique ID to the channel
	channel.ID = m.nextID
	m.nextID++

	// Store the channel in the repository
	m.Channels[channel.ID] = channel

	return &channel, nil
}

type MockUtil struct {
	GenerateAuthUrlFunc func(user model.User) string
	HandleCallbackFunc  func(r *http.Request) (*model.Channel, error)
}

func (m *MockUtil) GenerateAuthUrl(user model.User) string {
	if m.GenerateAuthUrlFunc != nil {
		return m.GenerateAuthUrlFunc(user)
	}
	return ""
}

func (m *MockUtil) HandleCallback(r *http.Request) (*model.Channel, error) {
	if m.HandleCallbackFunc != nil {
		return m.HandleCallbackFunc(r)
	}
	return nil, nil
}

func TestHandler_GenerateAuthLink(t *testing.T) {
	mockRepo := NewMockChannelRepository()
	mockConfig := config.Data{}
	mockUtil := MockUtil{}

	// create the handler
	handler := New(mockRepo, mockConfig, &mockUtil)

	// set the auth user and create the request
	user := &model.User{
		ID:    1,
		Email: "m@gmail.com",
	}

	req := httptest.NewRequest("GET", "/api/channel/auth/link", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	// create a response recorder
	rr := httptest.NewRecorder()

	// mock the return value of GenerateAuthUrl
	mockUtil.GenerateAuthUrlFunc = func(user model.User) string {
		return "http://google.auth.url.com"
	}

	// call the handler
	handler.GenerateAuthLink(rr, req)

	// check status
	assert.Equal(t, http.StatusOK, rr.Code)

	// check the response for the link
	var response url
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "http://google.auth.url.com", response.Plain)
}

func TestHandler_HandleCallback(t *testing.T) {
	mockRepo := NewMockChannelRepository()
	mockConfig := config.Data{}
	mockUtil := &MockUtil{}

	// Create the handler with mock dependencies
	handler := New(mockRepo, mockConfig, mockUtil)

	// create the request
	req := httptest.NewRequest("GET", "/api/google/authorize/callback", nil)

	// create response recorder
	rr := httptest.NewRecorder()

	mockUtil.HandleCallbackFunc = func(r *http.Request) (*model.Channel, error) {
		return &model.Channel{
			ID:        1,
			UserID:    1,
			YoutubeID: "test-id",
			Title:     "test-title",
		}, nil
	}

	handler.HandleCallback(rr, req)

	// check the status
	assert.Equal(t, http.StatusOK, rr.Code)

	// check channel added successfully
	assert.Len(t, mockRepo.Channels, 1)
	addedChannel, found := mockRepo.Channels[1]
	assert.True(t, found)
	assert.Equal(t, int64(1), addedChannel.ID)
	assert.Equal(t, int64(1), addedChannel.UserID)
}
