package beat

import (
	"bot/internal/commands"
	"bot/pkg/config"
	"bot/pkg/logging"
	tasks "bot/pkg/task"
	"time"
)

func eachBeatTasks() {

	log := logging.Log.With().
		Caller().
		Str("module", "beat").
		Logger()

	for _, b := range config.C.Beat {
		target, err := tasks.GetTarget(b.TargetName)
		if err != nil {
			log.Error().Err(err).Msg("get target")
		}
		task, err := tasks.GetTask(b.TaskName)
		if err != nil {
			log.Error().Err(err).Msg("get task")
		}
		ctx := commands.Context{
			TargetName: b.TargetName,
			Target:     target,
			Task:       task,
			Log:        log,
		}
		go commands.TaskCommands[task.Name](ctx)
	}

}

// 定期执行任务
func StartBeat() {

	for {
		eachBeatTasks()
		time.Sleep(3 * time.Second)
	}

}
