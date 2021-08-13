package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type employee struct {
	gorm.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *gorm.DB

var err error

var (
	emp = []employee{
		{Name: "sasa", Age: 12},
		{Name: "sa", Age: 21},
	}
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	var emps []employee
	db.Find(&emps)
	json.NewEncoder(w).Encode(&emps)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var emp employee
	db.First(&emp, params["id"])
	json.NewEncoder(w).Encode(&emp)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp employee
	_ = json.NewDecoder(r.Body).Decode(&emp)
	db.Create(&emp)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var emp employee
	db.First(&emp, params["id"])
	_ = json.NewDecoder(r.Body).Decode(&emp)
	db.Save(&emp)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var emp employee
	db.Delete(&emp, params["id"])
}

func main() {
	godotenv.Load()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("SSLMODE")
	password := os.Getenv("DB_PASSWORD")
	myRouter := mux.NewRouter()

	db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%v user=%s dbname=%s sslmode=%s password=%s", host, port, user, name, ssl, password))
	if err != nil {
		panic("Database access Failed!")
	}
	defer db.Close()

	db.AutoMigrate(&employee{})

	for index := range emp {
		db.Create(&emp[index])
	}

	myRouter.HandleFunc("/employees", GetEmployees).Methods("GET")
	myRouter.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
	myRouter.HandleFunc("/employees", CreateEmployee).Methods("POST")
	myRouter.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PUT")
	myRouter.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")

	handler := cors.Default().Handler(myRouter)

	log.Fatal(http.ListenAndServe(":4000", handler))
}
