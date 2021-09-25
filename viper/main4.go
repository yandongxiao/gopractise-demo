package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host Host
}
type Host struct {
	Address string `mapstructure:"address"`
	Port int `mapstructure:"port"`
}

func main () {
	// 读取配置文件
	viper.AddConfigPath(".")
	viper.SetConfigName("config4")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
}
