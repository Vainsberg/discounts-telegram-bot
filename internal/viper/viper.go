package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

var Pass, User string

func ViperUser() string {
	var err error

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	User = viper.GetString("UserbymySQL")

	return User

}

func ViperPass() string {
	var err error

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	Pass = viper.GetString("PassbymySQL")

	return Pass

}
