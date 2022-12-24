package bootstrap

import (
	"fmt"
	"gohub/pkg/config"
	"gohub/pkg/redis"
)

// SetupRedis初始化Redis
func SetupRedis() {
	//建立Redis连接
	redis.ConectRedis(fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"))
}
