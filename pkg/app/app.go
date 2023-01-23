package app

import (
	"gohub/pkg/config"
	"time"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}
func IsProduction() bool {
	return config.Get("app.env") == "production"
}
func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}

// URL传参path拼接站点的URL
func URL(path string) string {
	return config.Get("app.url") + path
}

// V1URL拼接带v1标识URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}
