package hash

import (
	"gohub/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用bcrypt 对密码进行加密
func BcryptHash(password string) string {
	// GenerateFromPassword 的第二个参数是cost值。建议大于12，数值越大耗费时间越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogIf(err)
	return string(bytes)
}

// BcryptCheck
//
//	@Description: 对比明文密码和数据库的哈希值
//	@param password
//	@param hash
//	@return bool
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// BcryptIsHashed
//
//	@Description: 判断字符串是否是哈希过的数据
//	@param str
//	@return bool
func BcryptIsHashed(str string) bool {
	//bcrypt加密后的长度等于60
	return len(str) == 60
}
