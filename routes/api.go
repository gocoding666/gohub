package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controller/api/v1/auth"
	"net/http"
)

func RegisterAPIRoutes(r *gin.Engine) {
	//测试一个v1的路由组，我们所有的v1版本的路由都将存放在这里
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			//判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			//判断 Email 是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			//发送验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)
		}
	}
	// 注册一个路由
	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World!",
		})
	})

}
