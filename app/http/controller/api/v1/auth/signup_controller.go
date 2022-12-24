// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controller/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"net/http"
)

// SignupController注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist 检测手机号是否注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	//请求对象
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupPhoneExist); !ok {
		return
	}
	//检查数据库并返回响应
	c.JSON(http.StatusOK, gin.H{"exist": user.IsPhoneExist(request.Phone)})
}

// IsEmailExist
//
//	@Description: 检测邮箱是否注册
//	@receiver sc
//	@param c
func (sc SignupController) IsEmailExist(c *gin.Context) {
	//初始化请求对象
	request := requests.SignupEmailExistRequest{}

	if ok := requests.Validate(c, &request, requests.SignupEmailExist); !ok {
		return
	}
	//验证数据库并返回响应
	c.JSON(http.StatusOK, gin.H{"exist": user.IsEmailExist(request.Email)})
}