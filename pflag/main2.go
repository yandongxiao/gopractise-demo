package main

import (
	"fmt"
	"github.com/spf13/pflag"
)

func main() {
	var (
		host string
		port int
	)

	// 结尾的Var表示支持将参数的值，绑定到变量
	pflag.StringVar(&host, "host", "127.0.0.1", "MySQL service host address.")
	// 结尾的P表示支持短选项
	pflag.IntVarP(&port, "port", "P", 3306, "MySQL service host port.")
	// 弃用 port 的简写形式
	pflag.CommandLine.MarkShorthandDeprecated("port", "please use --port only")

	// 弃用标志
	pflag.StringVar(&host, "Host", "127.0.0.1", "MySQL service host address.")
	pflag.CommandLine.MarkDeprecated("Host", "please use --host instead")

	// 如果没有指定命令行参数 host, 则 host 变量的值是 127.0.0.1
	// 如果指定了命令行参数 host，但是没有设置它的值，则 host 变量的值是 localhost
	// 通过 Lookup 方法，修改 pflag 对象的 field
	pflag.Lookup("host").NoOptDefVal = "localhost"

	pflag.Parse()

	fmt.Printf("host: %v\n", host)
	fmt.Printf("port: %v\n", port)

	// CommandLine 是默认的全局 FlagSet 对象
	fmt.Println(pflag.CommandLine.GetInt("port"))



	// 获取非命令行参数
	fmt.Printf("argument number is: %v\n", pflag.NArg())
	fmt.Printf("argument list is: %v\n", pflag.Args())
	fmt.Printf("the first argument is: %v\n", pflag.Arg(0))
}
