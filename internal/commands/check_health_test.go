package commands

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
)

func Test_onceCheckHealth(t *testing.T) {
	type args struct {
		bu      *healthUnit
		command string
		check   string
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
		}, "http://with_normal", "http://with_normal"}, 0},
		{"check failed", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_normal", "http://with_crash"}, 1},
		{"restart failed", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_crash", "http://with_timeout"}, 2},
		{"restarted", args{&healthUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, "http://with_normal", "http://with_timeout"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _ := checkHealth(tt.args.bu, tt.args.check, tt.args.command)
			if code != tt.code {
				error_msg := fmt.Sprintf("code error, real: %d, expect: %d", code, tt.code)
				panic(error_msg)
			}
		})
	}
}
