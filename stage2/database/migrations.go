package database

import (
	"fmt"

	"github.com/mryan-3/hng11/stage2/models"
	"gorm.io/gorm"
)

func MigrateDatabase(DB *gorm.DB) {
	fmt.Println("Running migration")

	DB.AutoMigrate(
        &models.User{},
	)

	fmt.Println("Migration ran!")
}
