package commands

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func Test_restart(t *testing.T) {
	type args struct {
		ru            *restartUnit
		ctx           context.Context
		check         string
		command       string
		checkInterval time.Duration
	}
	logger := log.With().
		Str("module", "restart").
		Logger()

	mockHttpClient := &MockHttpClient{}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"good", args{&restartUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, makeContext(), "http://with_normal", "http://with_normal", 0}, 0},
		{"command failed", args{&restartUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, makeContext(), "http://with_normal", "http://with_crash", 0}, 1},
		{"check failed", args{&restartUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, makeContext(), "http://with_crash", "http://with_normal", 0}, 2},
		{"without check", args{&restartUnit{
			logger,
			mockHttpClient,
			mockHttpClient,
		}, makeContext(), "", "http://with_normal", 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := restart(tt.args.ctx, tt.args.ru, tt.args.check, tt.args.command, tt.args.checkInterval); got != tt.want {
				t.Errorf("restart() = %v, want %v", got, tt.want)
			}
		})
	}
}
