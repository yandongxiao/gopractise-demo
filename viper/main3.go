package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
)

func main () {
	// 读取配置文件
	viper.AddConfigPath(".")	// // 把当前目录加入到配置文件的搜索路径中，可调用多次添加多个搜索路径
	viper.SetConfigName("config3")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func (e fsnotify.Event){
		fmt.Println("Config file changed: ", filepath.Base(e.Name))
		fmt.Println(viper.Get("host"))
	})

	select { }
}
