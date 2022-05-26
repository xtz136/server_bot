package commands

import (
	"context"
	"testing"
)

func Test_restart(t *testing.T) {
	type args struct {
		ros     *restartJob
		ctx     context.Context
		check   string
		command string
	}

	mockHttpClient := &MockHttpClient{}
	ros := &restartJob{
		checkClient:   mockHttpClient,
		commandClient: mockHttpClient,
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"good", args{ros, makeContext(), "http://with_normal", "http://with_normal"}, 0},
		{"command failed", args{ros, makeContext(), "http://with_normal", "http://with_crash"}, 1},
		{"check failed", args{ros, makeContext(), "http://with_crash", "http://with_normal"}, 2},
		{"without check", args{ros, makeContext(), "", "http://with_normal"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := restart(tt.args.ctx, tt.args.check, tt.args.command, tt.args.ros)
			if got != tt.want {
				t.Errorf("restart() = %v, want %v", got, tt.want)
			}
		})
	}
}
