package cmd

import (
	"gohub/bootstrap"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {
	//设置gin的运行模式,支持debug,release,test
	//release 会屏蔽调试信息，官方建议生产环境中使用
	//非release模式gin终端打印太多信息，干扰到我们程序中的Log
	//故此设置为release,有特殊情况手动改为debug即可
	gin.SetMode(gin.ReleaseMode)

	//gin实例
	router := gin.New()
	//初始化路由绑定
	bootstrap.SetupRoute(router)
	//运行服务器
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server,error:= " + err.Error())
	}
}
