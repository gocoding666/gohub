package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	btsConfig "gohub/config"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

func init() {
	// 加载config目录下的配置信息
	btsConfig.Initialize()
}
func main() {
	//配置初始化，依赖命令行 --env参数
	var env string
	flag.StringVar(&env, "env", "", "加载.env文件，如--env=testing 加载的是.env.testing文件")
	flag.Parse()
	config.InitConfig(env)
	//初始化Logger
	bootstrap.SetupLogger()
	// 设置gin的运行模式，支持debug,release,test
	//release会屏蔽调试信息，官方建议生产环境中使用
	//非release模型gin终端打印太多信息，干扰到我们程序中的Log
	// 故此设置为release,有特殊情况手动改为debug即可
	gin.SetMode(gin.ReleaseMode)
	// new 一个 Gin Engine 实例
	r := gin.New()
	//初始化DB
	bootstrap.SetupDB()
	bootstrap.SetupRedis()
	// 初始化路由绑定
	bootstrap.SetupRoute(r)

	//028654
	//E00uNvgBC7ETaREmhmnF
	//logger.Dump(captcha.NewCaptcha().VerifyCaptcha("E00uNvgBC7ETaREmhmnF", "028654"), "正确的答案")
	//verifycode.NewVerifyCode().SendSMS("17729732926")

	// 运行服务
	err := r.Run(":" + config.Get("app.port"))
	if err != nil {
		//错误处理，端口被占用了或者其他错误
		fmt.Printf("Gohub  start Error : %v ", err)
	}
	fmt.Printf("Gohub  started... ")
}
