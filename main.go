package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Entry struct {
	Date string `json:"date"`
	Day  string `json:"day"`
	Task string `json:"task"`
}

var db *sql.DB

func read_func(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Reading Initiated")
	res.Header().Set("Content-Type", "application/json")
	query := "SELECT * FROM entries"
	result, err := db.Query(query)
	defer result.Close()
	var output []Entry
	for result.Next() {
		var a, b, c string
		result.Scan(&a, &b, &c)
		fmt.Printf("%s %s %s\n", a, b, c)
		output = append(output, Entry{a, b, c})
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(res).Encode(output)
}

func create_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var new_entry Entry
	err := json.NewDecoder(req.Body).Decode(&new_entry)
	stmt, err := db.Prepare("INSERT INTO entries (date, day, task) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(new_entry.Date, new_entry.Day, new_entry.Task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(new_entry)
}

func update_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var updated_entry Entry
	err := json.NewDecoder(req.Body).Decode(&updated_entry)

	stmt, err := db.Prepare("UPDATE entries SET task = ? WHERE date = ?")
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	stmt.Exec(updated_entry.Task, updated_entry.Date)
	json.NewEncoder(res).Encode(updated_entry)
}

func delete_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var del_entry Entry
	err := json.NewDecoder(req.Body).Decode(&del_entry)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM entries WHERE date = ?")
	stmt.Exec(del_entry.Date)
	json.NewEncoder(res).Encode(del_entry)
}

func main() {
	router := mux.NewRouter()
	var err error
	defer db.Close()

	db, err = sql.Open("mysql", "root:@tcp(0.0.0.0:3306)/dbnew")
	err = db.Ping()
	//_, err = db.Exec("CREATE TABLE entries(date VARCHAR(20) PRIMARY KEY , day VARCHAR(20), task VARCHAR(100) );")

	if err != nil {
		log.Fatal(err)
	}
	router.HandleFunc("/", read_func).Methods("GET")
	router.HandleFunc("/", delete_func).Methods("DELETE")
	router.HandleFunc("/", create_func).Methods("POST")
	router.HandleFunc("/", update_func).Methods("PUT")

	fmt.Println("Starting server on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
