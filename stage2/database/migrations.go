package database

import (
	"fmt"

	"gorm.io/gorm"
)

func MigrateDatabase(DB *gorm.DB) {
	fmt.Println("Running migration")

	DB.AutoMigrate(
	)

	fmt.Println("Migration ran!")
}
