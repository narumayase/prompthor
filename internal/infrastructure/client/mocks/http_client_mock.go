package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of an HTTP client for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(url, contentType string, body []byte) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// CreateMockResponse is a helper function to create a mock HTTP response
func CreateMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}
