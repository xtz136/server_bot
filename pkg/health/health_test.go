package health

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
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

func Test_onceCheckHealth(t *testing.T) {
	type args struct {
		h       *healthUnit
		command string
		check   string
		wg      *sync.WaitGroup
	}
	logger := log.With().
		Caller().
		Str("module", "health").
		Logger()

	mockHttpClient := &MockHttpClient{}

	tests := []struct {
		name string
		args args
		code int
	}{
		{"good", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_normal", "http://with_normal", nil}, 0},
		{"restarted", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_normal", "http://with_timeout", nil}, 3},
		{"check failed", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_normal", "http://with_crash", nil}, 1},
		{"restart failed", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_crash", "http://with_timeout", nil}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			tt.args.wg = wg
			code, _ := onceCheckHealth(tt.args.h, tt.args.command, tt.args.check, tt.args.wg)
			if code != tt.code {
				error_msg := fmt.Sprintf("code error, real: %d, expect: %d", code, tt.code)
				panic(error_msg)
			}
		})
	}
}
