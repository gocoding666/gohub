package config

import (
	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper" //自定义包名，避免与内置viper实例冲突
	"gohub/pkg/helpers"
	"os"
)

// viper库实例
var viper *viperlib.Viper

// ConfigFunc动态加载配置信息
type ConfigFunc func() map[string]interface{}

// ConfigFuncs先加载到此数组，loadConfig再动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	// fmt.Println("配置模块-启动顺序验证-/pkg/config/config.go/init(){}")
	// 1.初始化viper库
	viper = viperlib.New()
	//2.配置类型，支撑“json”,"toml","yaml","yml","properties","props","prop","env","dotenv"
	viper.SetConfigType("env")
	// 3.环境变量配置文件查找的路径，相对于main.go
	viper.AddConfigPath(".")
	// 4.设置环境变量前缀，用以区分Go的系统环境变量
	viper.SetEnvPrefix("appenv")
	// 5.读取环境变量（支撑flags）
	viper.AutomaticEnv()
	ConfigFuncs = make(map[string]ConfigFunc)
}

// InitConfig初始化配置信息，完成对环境变量以及config信息的加载
func InitConfig(env string) {
	// fmt.Println("配置模块-启动顺序验证-/pkg/config/config.go/InitConfig(){}")
	//1.加载环境变量
	loadEnv(env)
	//2.注册配置信息
	loadConfig()
}

func loadEnv(envSuffix string) {
	// fmt.Println("配置模块-启动顺序验证-/pkg/config/config.go/loadEnv(){}")
	//默认加载.env文件，如果有传参--env=name的话，加载.env.name文件
	envPath := ".env"
	if len(envSuffix) > 0 {
		filepath := ".env." + envSuffix
		if _, err := os.Stat(filepath); err == nil {
			//如.env.testing或.env.stage
			envPath = filepath
		}
	}
	//加载env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//监控.env文件，变更时重新加载
	viper.WatchConfig()
}
func loadConfig() {
	// fmt.Println("配置模块-启动顺序验证-/pkg/config/config.go/loadConfig(){}")
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

// Env读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	// fmt.Println("配置模块-启动顺序验证-/pkg/config/config.go/Env(){}")
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

// Add新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// Get 获取配置项
// 第一个参数path允许使用点式获取，如app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}
func internalGet(path string, defaultValue ...interface{}) interface{} {
	//config或者环境变量不存在的情况
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString获取string类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt获取Int类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetFloat64获取float64类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetInt64获取int64类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint获取Uint类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool获取bool类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
