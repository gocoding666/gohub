// Package auth 授权相关逻辑
package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/logger"
)

func Attempt(email string, password string) (user.User, error) {
	userModel := user.GetByMulti(email)
	if userModel.ID == 0 {
		return user.User{}, errors.New("账号不存在")
	}
	if !userModel.ComparePassword(password) {
		return user.User{}, errors.New("密码错误")
	}
	return userModel, nil
}

func LoginByPhone(phone string) (user.User, error) {
	userModel := user.GetByPhone(phone)
	if userModel.ID == 0 {
		return user.User{}, errors.New("手机号未注册")
	}
	return userModel, nil
}

// CurrentUser从gin.Context中获取当前登录用户
func CurrentUser(c *gin.Context) user.User {
	userModel, ok := c.MustGet("current_user").(user.User)
	if !ok {
		logger.LogIf(errors.New("无法获取用户"))
		return user.User{}
	}
	//db is now a *DB value
	return userModel
}

// CurrentUID从gin.context中获取当前登录用户ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
