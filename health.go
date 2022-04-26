package main

import (
	"bot/internal/commands"
	"bot/pkg/logging"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

func iter_check_health() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log := logging.Log.With().
		Str("module", "health").
		Logger()

	checkClient := http.Client{
		Timeout: 1 * time.Second,
	}
	commadClient := http.Client{
		Timeout: 5 * time.Second,
	}
	var resp *http.Response
	var err error

	for {
		log.Info().Msg("start check health")
		time.Sleep(time.Duration(3) * time.Second)

		for key, _ := range viper.GetStringMap("systems") {
			b := commands.System{}
			viper.UnmarshalKey("systems."+key, &b)

			for _, item := range b.Health {
				command := item.Command
				check := item.Check

				resp, err = checkClient.Head(check)
				if err == nil {
					io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
					continue
				}

				if os.IsTimeout(err) {
					log.Info().Str("check", check).Msg("get bad")

					resp, err = commadClient.Get(command)
					if err != nil {
						log.Error().Err(err).Msg("restart")
						continue
					}
					log.Info().Str("service", command).Msg("restart")
					io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()

				} else {
					log.Error().Err(err).Msg("get error")
					continue
				}
			}
		}
	}

}
