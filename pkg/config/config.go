package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Variable struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type Task struct {
	Name    string `mapstructure:"name"`
	Command string `mapstructure:"command"`
	Check   string `mapstructure:"check"`
	Hidden  bool   `mapstructure:"hidden"`
}

type Target struct {
	Url []string `mapstructure:"url"`
}

type Beat struct {
	TargetName string `mapstructure:"target_name"`
	TaskName   string `mapstructure:"task_name"`
}

type Log struct {
	Level         int    `mapstructure:"level"`
	EnableConsole bool   `mapstructure:"enable_console"`
	EnableFile    bool   `mapstructure:"enable_file"`
	LogFileDir    string `mapstructure:"log_file_dir"`
	LogFileName   string `mapstructure:"log_file_name"`
}

type DingDing struct {
	Enable    bool   `mapstructure:"enable"`
	AppSecret string `mapstructure:"app_secret"`
}

type App struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	App       App               `mapstructure:"app"`
	Targets   map[string]Target `mapstructure:"targets"`
	Variables []Variable        `mapstructure:"variables"`
	Tasks     map[string]Task   `mapstructure:"tasks"`
	Beat      []Beat            `mapstructure:"beat"`
	Log       Log               `mapstructure:"log"`
	DingDing  DingDing          `mapstructure:"dingding"`
}

var C Config

func init() {
	viper.SetConfigName("task_bot")
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		path, _ := os.Getwd()
		fmt.Printf("fatal error config file %v in %v", err, path)
	}
	if err := viper.Unmarshal(&C); err != nil {
		panic(err)
	}
}
