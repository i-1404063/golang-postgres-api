package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct {
	sync.Mutex
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
	books  = []Book{
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

// Environment variables
var (
	host     = goDotEnv("HOST")
	dbPort   = goDotEnv("DBPORT")
	user     = goDotEnv("USER")
	dbName   = goDotEnv("NAME")
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

	// close connection when main function finishes
	// defer db.Close()

	// Make migrations to database
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	/// Insert row into the table
	// db.Create(person)
	// for idx:= range books {
	// 	db.Create(&books[idx])
	// }

	// route handler for person
	http.HandleFunc("/lib/users", getPeople)
	http.HandleFunc("/lib/user/", getPerson)
	http.HandleFunc("/lib/users/create", createPerson)
	http.HandleFunc("/lib/user/delete/", deletePerson)

	// route handler for book
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// all handler implemented here
func deletePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		path := strings.Split(r.URL.String(), "/")
		if len(path) != 5 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Provide valid id"))
			return
		}

		var person Person
		db.First(&person, path[4])
		db.Delete(&person)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)

	} else {
		w.Header().Set("Allowed method", http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only delete request is valid for this route"))
	}
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide proper information"))
			return
		}

		var person Person
		ct := r.Header.Get("content-type")
		if ct != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Header().Set("need content-type 'application/json' but got %s", ct)
			w.Write([]byte("Only Allowed Content-Type 'application/json'"))
			return
		}

		json.Unmarshal(bodyBytes, &person)

		createdPerson := db.Create(&person)
		err = createdPerson.Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something wrong happened"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)

	} else {
		w.Header().Set("Allowed method", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only post request is valid for this route"))
	}
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id is not provided"))
		return
	}

	var person Person
	var books []Book
	// db.Find(&person, "id = ?", path[3])
	db.First(&person, path[3])
	db.Preload("person").Find(&books, "person_id = ?", path[3]) // it will give all the books
	// db.Preload("person").First(&books, path[3]) // it will give only one book

	person.Books = books

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&person)
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	var person []Person

	db.Find(&person)
	json.NewEncoder(w).Encode(&person)
}
