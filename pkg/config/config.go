package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type State map[string]interface{}

type System struct {
	Name           string `mapstructure:"name"`
	LockIP         string `mapstructure:"lock_ip"`
	ListUserSystem string `mapstructure:"list_user_system"`
	MakeToken      string `mapstructure:"make_token"`
	Restart        []struct {
		Command string `mapstructure:"command"`
		Check   string `mapstructure:"check"`
	} `mapstructure:"restart"`
	Health []struct {
		Command string `mapstructure:"command"`
		Check   string `mapstructure:"check"`
	} `mapstructure:"health"`
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

type Config struct {
	Systems  map[string]System `mapstructure:"systems"`
	Log      Log               `mapstructure:"log"`
	DingDing DingDing          `mapstructure:"dingding"`
}

var C Config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Printf("fatal error config file %v in %v", err, path)
	}
	viper.Unmarshal(&C)
	fmt.Printf("config: %v\n", C)
}
