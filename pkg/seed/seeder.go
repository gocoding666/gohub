// Package seed 处理数据库填充相关逻辑
package seed

import "gorm.io/gorm"

//存放所有Seeder
var seeders []Seeder

//按顺序执行的Seeder数组
//支持一些必须按顺序执行的seeder,例如topic创建的
//时必须依赖与user,所以TopicSeeder 应该在UserSeeder后执行
var orderedSeederNames []string

type SeederFunc func(*gorm.DB)

//Seeder 对应每一个database/seeders 目录下的Seeder文件
type Seeder struct {
	Func SeederFunc
	Name string
}

//Add注册到seeders数组中
func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

//SetRunOrder 设置【按顺序执行的Seeder 数组】
func SetRunOrder(names []string) {
	orderedSeederNames = names
}
