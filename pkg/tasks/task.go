package tasks

import (
	"bot/pkg/config"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type TargetTask struct {
	Command string
	Check   string
}

func GetTask(taskName string) (*config.Task, error) {
	return getItem(config.C.Tasks, taskName)
}

func GetTarget(targetName string) ([]config.Target, error) {
	if v, ok := config.C.Targets[targetName]; ok {
		return v, nil
	} else {
		return nil, errors.New("item not found")
	}
}

func getItem[M any](m map[string]M, name string) (*M, error) {
	if v, ok := m[name]; ok {
		return &v, nil
	} else {
		return nil, errors.New("item not found")
	}
}

// 字符串模板替换，如果字符串中有{}，则替换为 p 中对应的值
func formatString(format string, p map[string]string) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func mapTargetTask(targets []config.Target, task *config.Task, vs map[string]string) []TargetTask {
	targetTasks := make([]TargetTask, len(targets))
	for i, target := range targets {
		serverUrl := target.Url
		urlVariables := map[string]string{}
		for k, v := range vs {
			urlVariables[k] = v
		}
		for _, v := range target.Variables {
			urlVariables[v.Name] = v.Value
		}
		baseUrl, _ := url.Parse(formatString(serverUrl, urlVariables))
		commandUrl, _ := baseUrl.Parse(formatString(task.Command, urlVariables))
		checkUrl, _ := baseUrl.Parse(formatString(task.Check, urlVariables))
		targetTasks[i] = TargetTask{
			Command: commandUrl.String(),
			Check:   checkUrl.String(),
		}
	}
	return targetTasks
}

// 列出指定的 target 和 task 组合，并根据 variables 替换
func ListTargetTask(targets []config.Target, task *config.Task, vs *[]config.Variable) []TargetTask {
	dict := make(map[string]string)
	for _, v := range *vs {
		dict[v.Name] = v.Value
	}
	return mapTargetTask(targets, task, dict)
}
