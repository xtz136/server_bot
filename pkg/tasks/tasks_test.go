package tasks

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakeJob struct {
	result int
}

func (r *fakeJob) Execute(ctx context.Context) int {
	result := r.GetResult()
	fmt.Printf("%v, job excute, result is %d\n", time.Now(), result)
	return result
}

func (r *fakeJob) SetResult(result int) {
	r.result = result
}

func (r *fakeJob) GetResult() int {
	return r.result
}

func makeConext() context.Context {
	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	return ctx
}

func makeJob(result int) *fakeJob {
	return &fakeJob{result}
}

func TestNewParallel(t *testing.T) {
	type args struct {
		ctx  context.Context
		jobs []Job
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"good", args{makeConext(), []Job{makeJob(0), makeJob(0), makeJob(0)}}, true},
		{"bad at first", args{makeConext(), []Job{makeJob(1), makeJob(0), makeJob(0)}}, false},
		{"bad at middle", args{makeConext(), []Job{makeJob(0), makeJob(1), makeJob(0)}}, false},
		{"bad at last", args{makeConext(), []Job{makeJob(0), makeJob(0), makeJob(1)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParallelTasks(tt.args.ctx, tt.args.jobs); got != tt.want {
				t.Errorf("NewParallelTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
