package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var (
	person = &Person{Name: "imon", Email: "imonhasans33@gmail.com"}
	books = []Book{
		{Title: "Captive Lady", Author: "Michael Madhusudhon Dutta", CallNumber: 25, PersonID: 1},
		{Title: "Hamlet", Author: "Shakespear", CallNumber: 20, PersonID: 2},
	}
)

var (
	db  *gorm.DB
	err error
)

func goDotEnv(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

//Environment variables
var (
	host = goDotEnv("HOST")
	dbPort = goDotEnv("DBPORT")
	user = goDotEnv("USER")
	dbName = goDotEnv("NAME")
	password = goDotEnv("PASSWORD")
)

func main() {
	// database connection string
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable port=%s password=%s TimeZone=Asia/Shanghai", host, user, dbName, dbPort, password)

	// opening connection to database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to the database!")
	}

	// Make migrations to database
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	/// Insert row into the table
	// db.Create(person)
	// for idx:= range books {
	// 	db.Create(&books[idx])
	// }
}
