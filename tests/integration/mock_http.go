package integration

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient implementa la interfaz HTTPClient para pruebas
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	args := m.Called(ctx, url)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHTTPClient) Post(ctx context.Context, url string, body []byte) ([]byte, error) {
	args := m.Called(ctx, url, body)
	return args.Get(0).([]byte), args.Error(1)
}
