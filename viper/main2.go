package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"strings"
)


func main () {
	var (
		hostSet string
		hostCommand string
		hostEnv string
		hostConfig string
		hostReader string
		hostDefault string
	)
	pflag.StringVar(&hostSet, "host-set", "127.0.0.1", "")
	pflag.StringVar(&hostCommand, "host-command", "127.0.0.1", "")
	pflag.StringVar(&hostEnv, "host-env", "127.0.0.1", "")
	pflag.StringVar(&hostConfig, "host-config", "127.0.0.1", "")
	pflag.StringVar(&hostReader, "host-reader", "127.0.0.1", "")
	pflag.StringVar(&hostDefault, "host-default", "127.0.0.1", "")
	pflag.Parse()

	// 设置默认值
	viper.Set("host-set", "host-set")

	// 绑定命令行参数
	viper.BindPFlags(pflag.CommandLine)

	// 绑定环境变量
	viper.AutomaticEnv()
	// viper.Get("a-b"), 那么
	// export 环境变量时，无法使用横杠，常见分隔符是下划线
	// viper.SetEnvKeyReplacer 使得 viper.Get("host-env") 对应的环境变量为 HOST_ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 读取配置文件
	viper.AddConfigPath(".")	// // 把当前目录加入到配置文件的搜索路径中，可调用多次添加多个搜索路径
	viper.SetConfigName("config2")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}


	// 从 io.Reader 获取配置
	viper.ReadConfig(strings.NewReader(`
host-set:     host-set-reader
host-command: host-command-reader
host-env: host-env-reader
host-reader: host-reader-reader
`))

	// 设置默认值
	viper.SetDefault("host-set", "host-set-default")
	viper.SetDefault("host-command", "host-command-default")
	viper.SetDefault("host-env", "host-env-default")
	viper.SetDefault("host-config", "host-config-default")
	viper.SetDefault("host-reader", "host-reader-default")
	viper.SetDefault("host-default", "host-default-default")

	fmt.Println(viper.Get("host-set"))
	fmt.Println(viper.Get("host-command"))
	fmt.Println(viper.Get("host-env"))
	fmt.Println(viper.Get("host-config"))
	fmt.Println(viper.Get("host-reader"))
	fmt.Println(viper.Get("host-default"))

	// 序列化
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	fmt.Println("===== 序列化 =====")
	fmt.Println(string(bs))
}