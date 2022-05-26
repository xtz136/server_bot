package tasks

import (
	"context"
	"sync"
)

type Job interface {
	Execute(ctx context.Context) int
	SetResult(int)
	GetResult() int
}

func NewParallelTasks(ctx context.Context, jobs []Job) bool {
	// 固定工作端的数量，避免机器CPU负荷太高
	max := 2
	queue := make(chan Job, max)
	output := make(chan Job)
	defer close(output)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 用来等待发起任务和接收任务两个协程都顺利完成
	wg := &sync.WaitGroup{}

	// 接收任务，并执行任务
	for i := 0; i < max; i++ {
		go func() {
			for job := range queue {
				select {
				case <-ctx.Done():
					wg.Done()
				default:
					job.SetResult(job.Execute(ctx))
					output <- job
					wg.Done()
				}
			}
		}()
	}

	// 发起任务
	go func() {
		defer close(queue)
		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				queue <- job
			}
		}
	}()

	// 当任务取消，清理现场
	tearDown := func() {
		// 把 output 还剩下的消息都抛弃掉
		go func() {
			for range output {
			}
		}()
		// 等待queue 和 output 关闭完成
		wg.Wait()
	}

	for i := 0; i < len(jobs); i++ {
		select {
		case resultjob := <-output:
			if resultjob.GetResult() != 0 {
				// 有一个任务出错了，剩下的就不需要做了
				cancel()
				<-ctx.Done()
				tearDown()
				return false
			}
		case <-ctx.Done():
			tearDown()
			return false
		}
	}
	return true
}
