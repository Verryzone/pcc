package main

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func koneksi() *gorm.DB {
	//membaca file .env
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	//user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
