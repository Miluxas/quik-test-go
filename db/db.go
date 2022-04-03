package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var db *gorm.DB

func Init() {
	log.Println("connecting")
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal(err)
	}

}
func connectDB() (*gorm.DB, error) {
	var err error
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp" + "(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}

func GetDB() *gorm.DB {
	return db
}
