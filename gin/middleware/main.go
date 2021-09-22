package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	content = ""
)

func hello(c *gin.Context) {
	log.Println("in hello")
	if content == "" {
		c.JSON(http.StatusOK, "ok")
	} else {
		c.JSON(http.StatusOK, content)
	}
	content = ""
}

// 你在 middleware 中没有调用c.Next()，默认这些操作在Handle之前执行
func middle(c *gin.Context) {
	if content == "" {
		content = "middleware"
	}
}

// 下面的操作在 middleware 之前执行
func auth(c *gin.Context) {
	if content == "" {
		content = "auth"
	}
}

// 常见middleware: https://static001.geekbang.org/resource/image/67/10/67137697a09d9f37bd87a81bf322f510.jpg?wh=1832x1521
func main() {
	// 创建一个不带任何中间件的路由
	r := gin.New()

	// 全局中间件
	// Logger 中间件将日志写到 gin.DefaultWriter，即使设置了 GIN_MODE=release
	// 默认设置 gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery 中间件，从任何 panic 恢复，并返回一个 500 错误
	r.Use(gin.Recovery())

	// 对于每一个路由，如果有需要，可以添加多个中间件
	r.GET("/benchmark", hello)

	// 授权组
	// authorized := r.Group("/", AuthRequired())
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"}))
	// 在这个示例中，我们使用了一个自定义的中间件 AuthRequired()，该中间件只作用于 authorized 组
	authorized.Use(middle)
	{
		authorized.POST("/hello", hello)

		// 嵌套组
		testing := authorized.Group("testing")
		testing.GET("/analytics", hello)
	}

	// 监听并服务于0.0.0.0:8080
	r.Run(":8080")
}
