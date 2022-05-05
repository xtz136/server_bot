package task

import (
	"bot/pkg/config"
	"reflect"
	"testing"
)

func Test_listTargetTask(t *testing.T) {
	type args struct {
		target *config.Target
		task   *config.Task
		vs     *[]config.Variable
	}
	tests := []struct {
		name string
		args args
		want []TargetTask
	}{
		{
			"good",
			args{
				&config.Target{Url: []string{"http://localhost:8080/api/{version}/"}},
				&config.Task{Command: "command?token={token}", Check: "check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
					{Name: "version", Value: "v1"},
				},
			},
			[]TargetTask{
				{
					"http://localhost:8080/api/v1/command?token=123",
					"http://localhost:8080/api/v1/check?token=123",
				},
			},
		},
		{
			"check not base target url",
			args{
				&config.Target{Url: []string{"http://localhost/api/v1/"}},
				&config.Task{Command: "http://other/command?token={token}", Check: "check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
				},
			},
			[]TargetTask{
				{
					"http://other/command?token=123",
					"http://localhost/api/v1/check?token=123",
				},
			},
		},
		{
			"task url contanin slash prefix",
			args{
				&config.Target{Url: []string{"http://localhost/api/v1/"}},
				&config.Task{Command: "/command?token={token}", Check: "check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
				},
			},
			[]TargetTask{
				{
					"http://localhost/command?token=123",
					"http://localhost/api/v1/check?token=123",
				},
			},
		},
		{
			"no variables",
			args{
				&config.Target{Url: []string{"http://localhost/api/v1/"}},
				&config.Task{Command: "command", Check: "check?token={token}"},
				&[]config.Variable{},
			},
			[]TargetTask{
				{
					"http://localhost/api/v1/command",
					"http://localhost/api/v1/check?token={token}",
				},
			},
		},
		{
			"more variables",
			args{
				&config.Target{Url: []string{"http://localhost/api/v1/"}},
				&config.Task{Command: "command", Check: "check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
					{Name: "token2", Value: "456"},
				},
			},
			[]TargetTask{
				{
					"http://localhost/api/v1/command",
					"http://localhost/api/v1/check?token=123",
				},
			},
		},
		{
			"target url is invalid url",
			args{
				&config.Target{Url: []string{"invalid_url"}},
				&config.Task{Command: "command", Check: "check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
				},
			},
			[]TargetTask{
				{
					"/command",
					"/check?token=123",
				},
			},
		},
		{
			"complex variables",
			args{
				&config.Target{Url: []string{"http://localhost/"}},
				&config.Task{Command: "command?var1={var1}", Check: "check?var2={var2}&var3={var3}"},
				&[]config.Variable{
					{Name: "var1", Value: "123 456"},
					{Name: "var2", Value: "[val]"},
					{Name: "var3", Value: "/val"},
				},
			},
			[]TargetTask{
				{
					"http://localhost/command?var1=123 456",
					"http://localhost/check?var2=[val]&var3=/val",
				},
			},
		},
		{
			"target url contains secondary directory",
			args{
				&config.Target{Url: []string{"http://localhost/website/"}},
				&config.Task{Command: "command?token={token}", Check: "../checker/check?token={token}"},
				&[]config.Variable{
					{Name: "token", Value: "123"},
				},
			},
			[]TargetTask{
				{
					"http://localhost/website/command?token=123",
					"http://localhost/checker/check?token=123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListTargetTask(tt.args.target, tt.args.task, tt.args.vs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listTargetTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
