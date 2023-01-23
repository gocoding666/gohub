package user

import (
	"gohub/pkg/app"
	"gohub/pkg/database"
	"gohub/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// IsEmailExist 判断Email已被注册
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// GetByPhone通过手机号来获取用户
func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone = ? ", phone).First(&userModel)
	return
}
func GetByMulti(loginID string) (userModel User) {
	database.DB.Where("phone = ? ", loginID).Or("email = ? ", loginID).Or("name= ?", loginID).First(&userModel)
	return
}

func Get(idstr string) (userModel User) {
	database.DB.Where("id", idstr).First(&userModel)
	return
}

func GetByEmail(email string) (userModel User) {
	database.DB.Where("email = ?", email).First(&userModel)
	return
}

// All 获取所有用户数据
func All() (users []User) {
	database.DB.Find(&users)
	return
}

func Paginate(c *gin.Context, perPage int) (users []User, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(User{}),
		&users,
		app.V1URL(database.TableName(&User{})),
		perPage,
	)
	return
}
