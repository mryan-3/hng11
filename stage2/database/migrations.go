package database

import (
	"fmt"

	"github.com/mryan-3/hng11/stage2/models"
	"gorm.io/gorm"
)

func MigrateDatabase(DB *gorm.DB) {
	fmt.Println("Running migration")

	DB.AutoMigrate(
		models.User{},
		models.Organisation{},
	)

    Session := DB.Session(&gorm.Session{PrepareStmt: true})
    if Session != nil {
        fmt.Println("success")
    }

	fmt.Println("Migration ran!")
}
