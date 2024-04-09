package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func getDbUrl() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return "postgresql://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
}

func InitDb() {
	dburl := getDbUrl()
	var err error

	// Connect to the database
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	log.Println("Database connected")

	// Enable uuid-ossp extension
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		panic("failed to enable uuid-ossp extension: " + err.Error())
	}
	log.Println("uuid-ossp extension enabled")

	// Migrate the schema
	err = DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledUrl{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	log.Println("Database migrated")
}

func GetDb() *gorm.DB {
	return DBConn
}
