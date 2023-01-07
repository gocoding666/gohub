package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controller/api/v1"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

// LoginController 用户控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {
	// 1.验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LogingByPhone); !ok {
		return
	}
	// 2.尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		response.Error(c, err, "账号不存在")
	} else {
		//登录成功
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}
