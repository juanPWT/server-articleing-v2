package config

import (
	"fmt"

	"gorm.io/gorm"
)

var db *gorm.DB = ConnectDB()

func Migrate(i interface{}) {

	db.AutoMigrate(&i)

	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&i)
	if err != nil {
		fmt.Printf("? failed migrate to DB: %v\n", err)
	} else {
		fmt.Printf("Migrate to DB success\n")
	}
}
