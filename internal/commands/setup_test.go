package commands

import (
	"context"
	"errors"
	"net/http"
)

type MockTimeoutError struct {
	error
}

func (e MockTimeoutError) Timeout() bool {
	return true
}

type MockHttpClient struct {
}

func (mhc *MockHttpClient) Fetch(req *http.Request) ([]byte, error) {
	if req.URL.String() == "http://with_timeout" {
		return nil, &MockTimeoutError{error: errors.New("timeout")}
	}
	if req.URL.String() == "http://with_crash" {
		return nil, errors.New("crash")
	}
	return nil, nil
}

func makeContext() context.Context {
	ctx := context.Background()
	return ctx
}
