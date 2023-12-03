package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

var Pass, User string

type Config struct {
	DbUser string
	DbPass string
}

func ViperUser() string {
	var err error
	var config Config

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	config.DbUser = viper.GetString("UserbymySQL")

	return User

}

func ViperPass() string {
	var err error
	var config Config

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	config.DbPass = viper.GetString("PassbymySQL")

	return Pass

}
