package main

import (
	"github.com/gin-gonic/gin"
	"github.com/herrluk/douyin/controller"
)

func main() {
	// 连接数据库
	db := controller.InitDB()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// gin创建默认路由器r
	r := gin.Default()
	initRouter(r)

	//r.SetTrustedProxies([]string{"192.168.50.85"})
	err := r.Run(":9999")
	if err != nil {
		if err != nil {
			panic("err:" + err.Error())
		}
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
