// Package validator 存放自定义规则和验证器
package validators

import (
	"errors"
	"fmt"
	"gohub/pkg/database"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/thedevsaddam/govalidator"
)

// 此方发会在初始化时执行，注册自定义表单验证规则
func init() {
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")
		//第一个参数，表名称，如users
		tableName := rng[0]
		//第二个参数，字段名称，如email或者phone
		dbFiled := rng[1]
		//第三个参数，排除ID
		var exceptID string
		if len(rng) > 2 {
			exceptID = rng[2]
		}
		//用户请求过来的数据
		requestValue := value.(string)
		//拼接SQL
		query := database.DB.Table(tableName).Where(dbFiled+" = ? ", requestValue)

		//如果传参第三个参数，加上SQL Where 过滤
		if len(exceptID) > 0 {
			query.Where("id != ? ", exceptID)
		}
		//查询数据库
		var count int64
		query.Count(&count)
		//验证不通过，数据库能找到对应的数据
		if count != 0 {
			//如果有自定义错误消息的话
			if message != "" {
				return errors.New(message)
			}
			//默认的错误消息
			return fmt.Errorf("%v 已被占用", requestValue)
		}
		//验证通过
		return nil

	})

	// max_cn:8 中文长度设定不超过8
	govalidator.AddCustomRule("max_cn", func(field, rule, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			//如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})
	//min_cn :2 中文长度设定不小于2
	govalidator.AddCustomRule("min_cn", func(field, rule, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长须需大于 %d 个字", l)
		}
		return nil
	})
}
