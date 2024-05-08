package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var db = make(map[string]string)

func read_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(db)
}

func create_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	db[message.Key] = message.Value
	json.NewEncoder(res).Encode(db[message.Key])
}

func update_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	db[message.Key] = message.Value
	json.NewEncoder(res).Encode(message)
}

func delete_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	db[message.Key] = message.Value
	delete(db, message.Key)
	json.NewEncoder(res).Encode(db)
}

func main() {
	router := mux.NewRouter()
	db["123"] = "456"
	db["343"] = "333"
	router.HandleFunc("/read", read_func).Methods("GET")
	router.HandleFunc("/delete", delete_func).Methods("DELETE")
	router.HandleFunc("/create", create_func).Methods("POST")
	router.HandleFunc("/update", update_func).Methods("PUT")

	fmt.Println("Starting server on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
