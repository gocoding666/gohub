package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/bootstrap"
)

func main() {
	// new 一个 Gin Engine 实例
	r := gin.New()
	// 初始化路由绑定
	bootstrap.SetupRoute(r)
	// 运行服务
	err := r.Run(":8000")
	if err != nil {
		//错误处理，端口被占用了或者其他错误
		fmt.Printf("Gohub  start Error : %v ", err)
	}
}
