package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbUser      string
	DbPass      string
	Apikey      string
	CountCron   string
	PlatiApiUrl string
}

func NewConfig() (*Config, error) {
	var err error

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	dbUser := viper.GetString("UserbymySQL")
	dbPass := viper.GetString("PassbymySQL")
	apiKey := viper.GetString("ApiKey")
	countCron := viper.GetString("CountCron")
	platiApiUrl := viper.GetString("PlatiApiUrl")

	return &Config{DbUser: dbUser, DbPass: dbPass, Apikey: apiKey, CountCron: countCron, PlatiApiUrl: platiApiUrl}, nil
}
