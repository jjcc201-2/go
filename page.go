package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Page struct {
	Title    string
	Filename string
	Content  string
}

type Pages []Page

var rawJson = []byte(`[{"Title":"First Page","Filename":"page1.txt","Content":"This is the 1st Page."},{"Title":"Second Page","Filename":"page2.txt","Content":"The 2nd Page is this."}]`)

func main() {
	// Decoding the JSON
	var pages Pages
	err := json.Unmarshal(rawJson, &pages)
	if err != nil {
		log.Fatal("Problem decoding JSON ", err)
	}

	fmt.Println(pages[0].Title)

	/*
		// Re-encode for demonstration purposes
		rejson, err := json.Marshal(pages)
		if err != nil {
			log.Fatal("Cannot encode to JSON ", err)
		}
		fmt.Fprintf(os.Stdout, "%s", rejson)
	*/

}

// Create email and add it to the user's outbox
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	var email Email

	if uuid, err := uuid.NewUUID(); err == nil { // If successful generating a new UUID for email
		fmt.Println(uuid)
		fmt.Println(r.Body)
		// w.Header().Set("Location", r.Host+"/outbox/user/"+uuid.String())
		if err := json.Unmarshal(r.Body, &email); err == nil {
			w.WriteHeader(http.StatusCreated)
			emails[user] = email
		} else {
			w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError) // Failed to generate UUID
	}
}
