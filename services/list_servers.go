package services

import (
	"github.com/spf13/viper"
)

func ListServices() []string {
	s := viper.GetStringMap("systems")
	servers := make([]string, 0, len(s))

	for k := range s {
		servers = append(servers, k)
	}

	return servers
}
