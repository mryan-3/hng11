package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

// Connect to Database
func ConnectDb() {
	dsn := os.Getenv("POSTGRES_URI")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect the database. \n", err)
	}

	log.Println("CONNECTED to the database")

	db.Logger = logger.Default.LogMode(logger.Info)

    MigrateDatabase(db)
	DB = Dbinstance{Db: db}
}
