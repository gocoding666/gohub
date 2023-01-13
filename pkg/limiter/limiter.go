// Package limiter 处理限流逻辑
package limiter

import (
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"strings"

	"github.com/gin-gonic/gin"
	limiterlib "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// GetKeyIP 获取limitor 的key,IP
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limitor的key,路由+IP,针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeTokeyString(c.FullPath()) + c.ClientIP()
}

func CheckRate(c *gin.Context, key string, formatted string) (limiterlib.Context, error) {
	var context limiterlib.Context
	rate, err := limiterlib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return context, err
	}
	//初始化存储，使用我们程序里共用的redis.Redis对象
	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiterlib.StoreOptions{
		//为limiter设置前缀，保持redis里数据的整洁
		Prefix: config.GetString("app.name") + ":limiter",
	})
	if err != nil {
		logger.LogIf(err)
		return context, err
	}
	//使用上面的初始化的 limiter.Rate对象和存储对象
	limiterObj := limiterlib.New(store, rate)

	//获取限流的结果
	if c.GetBool("limiter-once") {
		return limiterObj.Peek(c, key)
	} else {
		//确保多个路由组里调用LimitIP进行限流时，只增加一次访问次数
		c.Set("limiter-once", true)
		return limiterObj.Get(c, key)
	}

}

// routeTokeyString辅助方法，将URL中的/格式为-
func routeTokeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}
