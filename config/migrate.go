package config

import (
	"fmt"
	"server-article/model"

	"gorm.io/gorm"
)

var db *gorm.DB = ConnectDB()

func Migrate() {

	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Reset_password{}, &model.Category{}, &model.Article{}, &model.Body{}, &model.Comment{}, &model.Reply{})
	if err != nil {
		fmt.Printf("? failed migrate to DB: %v\n", err)
	} else {
		fmt.Println("Migrate to DB success")
	}
}
