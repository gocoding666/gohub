package migrations

import (
	"database/sql"
	"gohub/app/models"
	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type User struct {
		models.BaseModel
	}
	type Category struct {
		models.BaseModel
	}

	type Topic struct {
		models.BaseModel
		Title      string `gorm:"type:varchar(255);not null;index"`
		Body       string `gorm:"type:longtext;not null"`
		UserID     string `gorm:"type:bigint;not null;index"`
		CategoryID string `gorm:"type:bigint;not null;index"`
		//会创建user_id和category_id 外键约束
		User     User
		Category Category

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Topic{})
	}
	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Topic{})
	}
	migrate.Add("2023_01_25_211350_add_topics_table", up, down)
}
