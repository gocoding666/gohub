package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controller/api/v1"
	"gohub/pkg/captcha"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

// VerifyCOdeController用户控制器
type VerifyCodeController struct {
	v1.BaseAPIController
}

func (v *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	//生产验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	//记录错误日志，因为验证码是用户的入口，出错时应该记error等级的日志
	logger.LogIf(err)
	//返回给用户
	response.JSON(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}
