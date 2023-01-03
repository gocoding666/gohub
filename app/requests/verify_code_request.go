package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/captcha"
)

type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

// VerifyCodePhone验证表单，返回长度等于零即通过
func VerifyCodePhone(data interface{}, c *gin.Context) map[string][]string {
	//1.定制认证规则
	rules := govalidator.MapData{
		"phone":          []string{"required", "digits:11"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}
	//2.定制错误消息
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称phone",
			"digits:手机号长度必须为11位的数字",
		},
		"captcha_id": []string{
			"required:图片验证码的ID为必填"},
		"captcha_answer": []string{
			"required:图片验证码答案必填",
			"digits:图片验证码长度必须为6位的数字"},
	}
	errs := validate(data, rules, messages)
	//图片验证码
	_data := data.(*VerifyCodePhoneRequest)
	if ok := captcha.NewCaptcha().VerifyCaptcha(_data.CaptchaID, _data.CaptchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "图片验证码错误")
	}
	return errs
}
