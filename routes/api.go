package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controller/api/v1/auth"
	"gohub/app/http/middlewares"
	authpkg "gohub/pkg/auth"
	"gohub/pkg/response"
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
			// 使用手机和验证码进行注册
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)
			//发送验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)

			lgc := new(auth.LoginController)
			//使用手机号，短信验证码进行登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			//支持手机号，Email和用户名
			authGroup.POST("/login/using-password", lgc.LoginByPassword)

			authGroup.POST("/login/refresh-token", lgc.RefreshToken)

			//重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", pwc.ResetByEmail)
		}
	}
	v1.GET("/test_auth", middlewares.AuthJWT(), func(c *gin.Context) {
		userModel := authpkg.CurrentUser(c)
		response.Data(c, userModel)

	})
	v1.GET("/test_guest", middlewares.GuestJWT(), func(c *gin.Context) {
		c.String(http.StatusOK, "Hello guest")
	})
}
