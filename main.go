package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var schema = `CREATE SEQUENCE IF NOT EXISTS public.employees_id_seq
INCREMENT 1
START 1
MINVALUE 1
MAXVALUE 2147483647
CACHE 1;
CREATE TABLE IF NOT EXISTS employees
 (
	id integer NOT NULL DEFAULT nextval('employees_id_seq'::regclass),
	 name text COLLATE pg_catalog."default",
	 age integer,
	 CONSTRAINT employees_pkey PRIMARY KEY (id)
 )
 `

type employee struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
	Age  int    `db:"age" json:"age,omitempty"`
}

var err error

func Connect() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("SSLMODE")
	password := os.Getenv("DB_PASSWORD")
	db := sqlx.MustOpen("postgres", fmt.Sprintf("host=%s port=%v user=%s dbname=%s sslmode=%s password=%s", host, port, user, name, ssl, password))
	return db

}
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	var emps []employee
	err = db.Select(&emps, "SELECT ID,Name,Age FROM employees")
	if err != nil {
		fmt.Fprintln(w, "Could not retreieve Employees")
		return
	}
	json.NewEncoder(w).Encode(emps)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	params := mux.Vars(r)
	var emp employee
	err = db.Get(&emp, "SELECT id,name,age from employees WHERE id=$1", params["id"])
	if err != nil {
		fmt.Fprintln(w, "Could not retreieve Employee")
		return
	}
	json.NewEncoder(w).Encode(&emp)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	var emp employee
	_ = json.NewDecoder(r.Body).Decode(&emp)
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO employees (Name, Age) VALUES ($1, $2)", emp.Name, emp.Age)
	tx.Commit()

	fmt.Fprint(w, "Creation Successful")
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	params := mux.Vars(r)
	var emp employee
	_ = json.NewDecoder(r.Body).Decode(&emp)
	tx := db.MustBegin()
	tx.MustExec("UPDATE employees SET Name=$1, Age=$2 WHERE id=$3", emp.Name, emp.Age, params["id"])
	tx.Commit()

	fmt.Fprint(w, "Updation Successful")
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	params := mux.Vars(r)
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM employees WHERE id=$1", params["id"])
	tx.Commit()

	fmt.Fprint(w, "Deletion Successful")
}

func main() {
	godotenv.Load()

	db := Connect()

	myRouter := mux.NewRouter()

	db.MustExec(schema)

	myRouter.HandleFunc("/employees", GetEmployees).Methods("GET")
	myRouter.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
	myRouter.HandleFunc("/employees", CreateEmployee).Methods("POST")
	myRouter.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PUT")
	myRouter.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")

	handler := cors.Default().Handler(myRouter)

	fmt.Println("Server up...")
	log.Fatal(http.ListenAndServe(":4001", handler))
}
