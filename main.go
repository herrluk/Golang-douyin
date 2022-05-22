package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	initDB()

	// gin创建默认路由器r
	r := gin.Default()

	initRouter(r)

	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
