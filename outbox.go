package main

import (
	// Use JSON
	"encoding/json"
	"fmt"
	"log"      // Simple logger
	"net/http" // Provides HTTP client and server implementations. GET, POST, HEAD and PostForm

	"github.com/google/uuid"
	"github.com/gorilla/mux" // A request router and dispatcher for matching incoming requests against a list of registered routes
)

type Email struct { // Type is used to refer to the struct afterwards
	from    string
	to      string
	message string
}

var emails map[string]Email // Used as a wimpy way of "storing" replies on the server

/*
1. Add email to user's outbox
2. List emails
*/

// List emails from a user's outbox
func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	fmt.Println("aaaaaaaaaaa: ", emails[user])

	if email, ok := emails[user]; ok {
		fmt.Println("EMAIL VALUE: ", email)
		w.WriteHeader(http.StatusOK)
		if enc, err := json.Marshal(email.message); err == nil { // If you have an error converting it to JSON
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

	if uuid, err := uuid.NewUUID(); err == nil { // If successful generating a new UUID for email
		fmt.Println(uuid)
		// w.Header().Set("Location", r.Host+"/outbox/user/"+uuid.String())
		if err := decoder.Decode(&email); err == nil { //If no errors in decoding the message into the email object
			w.WriteHeader(http.StatusCreated)
			emails[user] = email
			fmt.Println("abc ", email)
		} else {
			w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError) // Failed to generate UUID
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true) //
	router.HandleFunc("/outbox/{user}", Read).Methods("GET")
	router.HandleFunc("/outbox/{user}", Create).Methods("POST")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func main() {
	emails = make(map[string]Email)
	handleRequests()
}
