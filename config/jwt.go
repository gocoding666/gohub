package config

import "gohub/pkg/config"

func init() {
	config.Add("jwt", func() map[string]interface{} {
		return map[string]interface{}{
			//使用config.GetString("app.key")
			//"signing_key":
			"expire_time": config.Env("JWT_EXPIRE_TIME", 120),
			//允许刷新时间，单位分钟，86400为两个月，从Token的签名时间算起
			"max_refresh_time": config.Env("JWT_MAX_REFRESH_TIME", 86400),
			//debug模式下的过期时间，方便本地开发调试
			"debug_expire_time": 86400,
		}
	})
}
