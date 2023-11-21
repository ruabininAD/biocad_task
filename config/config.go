package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Init(ConfigName string) {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/configuration")
	viper.AddConfigPath("C:\\Users\\ryabi\\GolandProjects\\biocadGo\\config\\configuration")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}
