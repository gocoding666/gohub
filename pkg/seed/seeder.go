// Package seed 处理数据库填充相关逻辑
package seed

import (
	"gohub/pkg/console"
	"gohub/pkg/database"

	"gorm.io/gorm"
)

// 存放所有Seeder
var seeders []Seeder

// 按顺序执行的Seeder数组
// 支持一些必须按顺序执行的seeder,例如topic创建的
// 时必须依赖与user,所以TopicSeeder 应该在UserSeeder后执行
var orderedSeederNames []string

type SeederFunc func(*gorm.DB)

// Seeder 对应每一个database/seeders 目录下的Seeder文件
type Seeder struct {
	Func SeederFunc
	Name string
}

// Add注册到seeders数组中
func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

// SetRunOrder 设置【按顺序执行的Seeder 数组】
func SetRunOrder(names []string) {
	orderedSeederNames = names
}

// GetSeeder 通过名称来获取Seeder对象
func GetSeeder(name string) Seeder {
	for _, sdr := range seeders {
		if name == sdr.Name {
			return sdr
		}
	}
	return Seeder{}
}

// RunAll运行所有Seeder
func RunAll() {
	//先运行orderd的
	executed := make(map[string]string)
	for _, name := range orderedSeederNames {
		sdr := GetSeeder(name)
		if len(sdr.Name) > 0 {
			console.Warning("Running Ordered Seeder : " + sdr.Name)
			sdr.Func(database.DB)
			executed[name] = name
		}
	}

	//在运行升下的
	for _, sdr := range seeders {
		//过滤已运行的
		if _, ok := executed[sdr.Name]; !ok {
			console.Warning("Running Seeder : " + sdr.Name)
			sdr.Func(database.DB)
		}
	}
}

// RunSeeder 运行单个Seeder
func RunSeeder(name string) {
	for _, sdr := range seeders {
		if name == sdr.Name {
			sdr.Func(database.DB)
			break
		}
	}
}
