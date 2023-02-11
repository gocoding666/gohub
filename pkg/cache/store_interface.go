package cache

import "time"

type Store interface {
	Set(key string, value string, expireTime time.Duration)
	Get(key string) string
	Has(key string) bool
	Forget(key string)
	Forever(key string, value string)
	Flush()

	IsAlive() error

	//Increment当参数只有1个时，为key,增加1.
	//当参数有2个时，第一个参数为key,第二个参数为要增加的值 int64 类型。
	Increment(parameters ...interface{})
	//Decrement 当参数只有1个时，为key,减去1.
	//当参数有2个时，第一个参数为key,第二个参数为要减去的值 int64类型。
	Decrement(parameters ...interface{})
}
