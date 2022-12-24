package captcha

import (
	"github.com/mojocn/base64Captcha"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"sync"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// once 确保internalCaptcha对象只初始化一次
var once sync.Once

// internalCaptcha 内部使用的Captcha对象
var internalCaptcha *Captcha

// NewCaptcha 单例模式获取
func NewCaptcha() *Captcha {
	once.Do(func() {
		//初始化Captcha对象
		internalCaptcha = &Captcha{}
		//使用全局Redis对象，并配置存储Key的前缀
		store := RedisStore{
			RedisClient: redis.Redis,
			KeyPrefix:   config.GetString("app.name") + ":captcha:",
		}
		//配置base64Captcha驱动信息
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("captcha.height"),      //宽
			config.GetInt("captcha.width"),       //高
			config.GetInt("captcha.length"),      //长度
			config.GetFloat64("captcha.maxskew"), //数字的最大斜倾角度
			config.GetInt("captcha.dotcount"),    //图片背景里的混淆点数量
		)
		//实例化base64Captcha 并赋值给内部使用的internalCaptcha对象
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})
	return internalCaptcha
}

// GenerateCaptcha
//
//	@Description: 生成图片验证码
//	@receiver c
//	@return id
//	@return b64s
//	@return err
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}
func (c *Captcha) VerifyCaptcha(id string, answer string) (match bool) {
	//方便本地和API自动测试
	if !app.IsProduction() && id == config.GetString("captcha.testing_key") {
		return true
	}
	//第三个参数时验证后是否删除，我们选择false
	//这样方便用户多次，防止表单提交错误需要多次输入图片验证码
	return c.Base64Captcha.Verify(id, answer, false)
}
