package response

import (
	"github.com/gin-gonic/gin"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

// JSON响应200和JSON数据
func JSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Success响应200和预设【操作成功！】的JSON数据
// 执行某个【没有具体返回数据】的【变更】操作成功后调用，例如删除、修改密码、修改手机号
func Success(c *gin.Context) {
	JSON(c, gin.H{
		"success": true,
		"message": "操作成功！",
	})
}

// Data 响应200 和带data键的JSON数据
// 执行【更新操作】成功后调用，例如更新话题，成功后返回已更新的话题
func Data(c *gin.Context, data interface{}) {
	JSON(c, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created响应201和带data键的JSON数据
// 执行【更新操作】成功后调用，例如更新话题，成功后返回已更新的话题
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// CreatedJSON响应201和JSON数据
func CreatedJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}
func Abort404(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"message": defaultMessage("数据不存在，请确认请求正确", msg...),
	})
}
func Abort403(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": defaultMessage("权限不足，请确定您有对应的权限", msg...),
	})
}
func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": defaultMessage("服务器内部错误，请稍后再试", msg...),
	})
}

// BasRequest响应400，传参err对象，未传参msg时使用默认消息
// 在解析用户请求，请求和格式和方法不符合预期时调用
func BadRequest(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用multipart标头，参数请使用JSON格式。", msg...),
		"error":   err.Error(),
	})
}
func Error(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	//error类型未【数据库未找到内容】
	if err == gorm.ErrRecordNotFound {
		Abort404(c)
		return
	}
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": defaultMessage("请求处理失败，请查看error的值"),
		"error":   err.Error(),
	})
}

// Unauthorized 响应401，未传参msg时使用默认消息
// 登录失败、JWT解析失败时调用
func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": defaultMessage("请求解析失败，请确认请求格式是否正确。上传文件请使用multipart标头，参数请使用JSON格式。", msg...),
	})
}

func ValidationError(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "请求验证不通过，具体请查看errors",
		"errors":  errors,
	})
}

func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}
