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

// ConnectDb connects to the main database
func ConnectDb() {
    connectToDb(os.Getenv("POSTGRES_URI"))
}

// ConnectTestDb connects to the test database
func ConnectTestDb() {
    connectToDb(os.Getenv("TEST_POSTGRES_URI"))
}

// connectToDb is a helper function to connect to a database
func connectToDb(dsn string) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        PrepareStmt: false,
    })
    if err != nil {
        log.Fatal("Failed to connect to the database. \n", err)
    }
    log.Println("CONNECTED to the database")
    db.Logger = logger.Default.LogMode(logger.Info)
    MigrateDatabase(db)
    DB = Dbinstance{Db: db}
}
