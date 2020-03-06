package main

import (
	// Use JSON
	"encoding/json"
	"log"      // Simple logger
	"net/http" // Provides HTTP client and server implementations. GET, POST, HEAD and PostForm

	"github.com/gorilla/mux" // A request router and dispatcher for matching incoming requests against a list of registered routes
)

type Email struct { // Type is used to refer to the struct afterwards
	From    string
	To      string
	Message string
}

type Emails []Email

// type Emails map[string]Email

var emails map[string]Email // ONE EMAIL PER USER

var outbox map[string]Emails

/*
1. Add email to user's outbox
2. List emails
*/

// List emails from a user's outbox
func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	if emails, ok := outbox[user]; ok {

		// Obtain all the messages in the outbox of a specific user
		var emailMessages = make(map[int]string)
		for key, value := range emails {
			emailMessages[key] = value.Message
		}
		w.WriteHeader(http.StatusOK)
		if enc, err := json.Marshal(emailMessages); err == nil { // If you have an error converting it to JSON
			w.Write([]byte(enc))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Create email and add it to the user's outbox
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	decoder := json.NewDecoder(r.Body)
	var email Email

	if err := decoder.Decode(&email); err == nil { //If no errors in decoding the message into the email object
		w.WriteHeader(http.StatusCreated)
		outbox[user] = append(outbox[user], email)
	} else {
		w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
	}
}

// Read a particular email from the outbox
func Read(w http.ResponseWriter, r *http.Request) {

}

// Delete a particular email from the outbox
func Delete(w http.ResponseWriter, r *http.Request) {

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/outbox/{user}/{uuid}", Delete).Methods("DELETE")
	router.HandleFunc("/outbox/{user}/{uuid}", Read).Methods("GET")
	router.HandleFunc("/outbox/{user}", List).Methods("GET")
	router.HandleFunc("/outbox/{user}", Create).Methods("POST")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func main() {
	outbox = make(map[string]Emails)
	handleRequests()
}
