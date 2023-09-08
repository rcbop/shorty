package shortener

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// MockRedisClient is a mock implementation of a Redis client for testing purposes.
type MockRedisClient struct {
	existsResult int64
	setResult    error
	getResult    string
}

func (m *MockRedisClient) Exists(_ context.Context, key string) (int64, error) {
	return m.existsResult, nil
}

func (m *MockRedisClient) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	return m.setResult
}

func (m *MockRedisClient) Get(_ context.Context, key string) (string, error) {
	return m.getResult, nil
}

func TestRedisURLShortener_ShortenURL(t *testing.T) {
	// Create a mock Redis client for testing
	mockClient := &MockRedisClient{}

	// Create a RedisURLShortener instance with the mock client
	urlShortener := NewRedisURLShortener(mockClient)

	// Test case 1: Successfully shorten the URL
	mockClient.existsResult = 0
	mockClient.setResult = nil

	shortURL, err := urlShortener.ShortenURL("https://example.com")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(shortURL) != 6 {
		t.Errorf("Expected short URL length of 6, but got %d", len(shortURL))
	}

	// Test case 2: Redis operation fails
	mockClient.existsResult = 0
	mockClient.setResult = redis.ErrClosed

	_, err = urlShortener.ShortenURL("https://example.com")
	if err != redis.ErrClosed {
		t.Errorf("Expected error %v, but got %v", redis.ErrClosed, err)
	}

	// Add more test cases as needed
}

type ErrorMockRedisClient struct {
	shouldErrorExists bool
	shouldErrorSet    bool
}

func (m *ErrorMockRedisClient) Exists(ctx context.Context, key string) (int64, error) {
	if m.shouldErrorExists {
		return 0, errors.New("mock Redis error")
	}
	return 0, nil
}

func (m *ErrorMockRedisClient) Set(ctx context.Context, key string, value interface{}, _ time.Duration) error {
	if m.shouldErrorSet {
		return errors.New("mock Redis error")
	}
	return nil
}

func (m *ErrorMockRedisClient) Get(ctx context.Context, key string) (string, error) {
	return key, nil
}

func TestRedisURLShortener_ShortenURL_ErrorExists(t *testing.T) {
	tc := map[string]RedisClient{
		"testErrorExists": &ErrorMockRedisClient{
			shouldErrorExists: true,
			shouldErrorSet:    false,
		},
		"testErrorSet": &ErrorMockRedisClient{
			shouldErrorExists: false,
			shouldErrorSet:    true,
		},
	}
	for tt, v := range tc {
		t.Run(tt, func(t *testing.T) {
			urlShortener := NewRedisURLShortener(v)

			longURL := "https://example.com"
			_, err := urlShortener.ShortenURL(longURL)
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			expectedErrMsg := "mock Redis error"
			if err.Error() != expectedErrMsg {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
			}
		})
	}
}
