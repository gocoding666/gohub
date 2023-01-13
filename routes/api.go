package routes

import (
	"gohub/app/http/controller/api/v1/auth"
	"gohub/app/http/middlewares"
	authpkg "gohub/pkg/auth"
	"gohub/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	//测试一个v1的路由组，我们所有的v1版本的路由都将存放在这里
	v1 := r.Group("/v1")
	//全局限流中间件，每小时限流，这里是所有API(根据IP)请求加起来
	//作为参考Github API每小时最多60个请求（根据IP)
	//测试时，可以调高一点
	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		//限流中间件：每小时限流，作为参考Github API每小时最多60个请求（根据IP）
		//测试时，可以调高一点
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			//登录
			lgc := new(auth.LoginController)
			authGroup.POST("/login/using-phone", middlewares.GuestJWT(), lgc.LoginByPhone)
			authGroup.POST("/login/using-password", middlewares.GuestJWT(), lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			//重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), pwc.ResetByEmail)

			//注册用户
			suc := new(auth.SignupController)
			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsEmailExist)

			//发送验证码
			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)

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
