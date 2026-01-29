package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockConnector_Start(t *testing.T) {
	mock := NewMockConnector()
	assert.NoError(t, mock.Start(nil))
}

func TestMockConnector_Stop(t *testing.T) {
	mock := NewMockConnector()
	assert.NoError(t, mock.Stop())
}

func TestMockConnector_Name(t *testing.T) {
	mock := NewMockConnector()
	assert.Equal(t, "MockConnector", mock.Name())
}

func TestMockConnector_SendMessage(t *testing.T) {
	mock := NewMockConnector()
	assert.NoError(t, mock.SendMessage(nil, "user1", "test message"))
}

func TestMockConnector_Events(t *testing.T) {
	mock := NewMockConnector()
	events := mock.Events()
	assert.NotNil(t, events)
}
