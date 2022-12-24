package config

import "gohub/pkg/config"

func init() {
	config.Add("log", func() map[string]interface{} {
		return map[string]interface{}{
			"level": config.Env("LOG_LEVEL", "debug"),
			//日志的类型，可选：
			// “single” 独立的文件
			// “daily” 按照日期每日一个
			"type": config.Env("LOG_TYPE", "single"),
			//滚动日志配置
			//日志文件路径
			"filename": config.Env("LOG_NAME", "storage/logs/logs.log"),
			// 每个日志文件保存的最大尺寸 单位：M
			"max_size": config.Env("LOG_MAX_SIZE", 64),
			//最大保存日志文件数，0为不限，MaxAge到了还是会删除
			"max_backup": config.Env("LOG_MAX_BACKUP", 5),
			// 最多保存多少天，7表示一周前的日志会被删除，0表示不删
			"max_age": config.Env("LOG_MAX_AGE", 30),
			// 是否压缩，压缩日志不方便查看，我们设置为false
			"compress": config.Env("LOG_COMPRESS", false),
		}
	})
}
