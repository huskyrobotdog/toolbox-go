package config

import (
	"strings"

	"github.com/huskyrobotdog/toolbox-go/inner"
	"github.com/spf13/viper"
)

func Initialization[C any](path string) *C {
	c := new(C)
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "-"))
	if err := viper.ReadInConfig(); err != nil {
		inner.Panic(err.Error())
	}
	if err := viper.Unmarshal(c); err != nil {
		inner.Panic(err.Error())
	}
	return c
}
