package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

// ValidatorFunc验证函数类型
type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {
	//1.解析请求，支持Json数据、表单请求和URL Query
	if err := c.ShouldBindJSON(obj); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "请求解析错误，请确认请求格式是否正确。上传文件请使用multipart标头，参数请使用JSON格式。",
			"error":   err.Error(),
		})
		fmt.Println(err.Error())
		return false
	}
	//表单验证
	errs := handler(obj, c)
	//3.判断验证是否通过
	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "请求验证不通过，具体请查看errors",
			"errors":  errs,
		})
		return false
	}
	return true
}

func validate(data interface{}, rule govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Data:          data,
		Rules:         rule,
		TagIdentifier: "valid",
		Messages:      messages,
	}
	//开始验证
	return govalidator.New(opts).ValidateStruct()
}
