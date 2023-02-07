package factories

import (
	"gohub/app/models/link"

	"github.com/bxcodec/faker/v3"
)

func MakeLinks(times int) []link.Link {
	var objs []link.Link

	//设置唯一性，如Link 模型的某个字段需要唯一，即可取消注释
	//faker.SetGenerateUniqueValues(true)
	for i := 0; i < times; i++ {
		model := link.Link{
			Name: faker.Username(),
			URL:  faker.URL(),
		}
		objs = append(objs, model)

	}
	return objs
}
