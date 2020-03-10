package main

import (
	// Use JSON
	"encoding/json"
	"log"      // Simple logger
	"net/http" // Provides HTTP client and server implementations. GET, POST, HEAD and PostForm

	"github.com/google/uuid"
	"github.com/gorilla/mux" // A request router and dispatcher for matching incoming requests against a list of registered routes
)

// An email is made out of these constituents
type Email struct {
	From    string
	To      string
	Message string
}

// Each user has an inbox and an outbox.
type User struct {
	Inbox  map[string]Email // An inbox is a map of all emails inside it
	Outbox map[string]Email // An outbox is also a map of all emails inside it
}

var mailSubmissionAgent map[string]User

// List emails in a user's inbox or outbox
func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]

	if account, ok := mailSubmissionAgent[user]; ok { // Does the user exist?

		if box == "inbox" || box == "outbox" { // Has the box been specified?

			var chosenbox map[string]Email // Removes the necesity to duplicate fucntions for both the inbox and outbox

			if box == "inbox" {
				chosenbox = account.Inbox
			} else {
				chosenbox = account.Outbox
			}

			// Obtain all the messages
			var emailMessages = make(map[string]string)
			for emailUuid, emailBody := range chosenbox {
				emailMessages[emailUuid] = emailBody.Message
			}
			if enc, err := json.Marshal(emailMessages); err == nil { // If you have an error converting it to JSON
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(enc))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest) // Must state whether targeting the inbox or outbox
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // User does not exist
	}
}

// Create email and add it to the user's outbox
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	decoder := json.NewDecoder(r.Body)
	var email Email

	// UUID to uniquely identify an email
	if uuid, err := uuid.NewUUID(); err == nil {
		if err := decoder.Decode(&email); err == nil { //If no errors in decoding the message into the email object

			// If user has not been created, create them
			if mailSubmissionAgent[user].Outbox == nil { // PROBLEM: if you instaniate a new inbox and outbox based on the fact that the outbox is empty, their inbox could be rewritten
				mailSubmissionAgent[user] = User{Inbox: make(map[string]Email), Outbox: make(map[string]Email)}
			}

			w.WriteHeader(http.StatusCreated)
			mailSubmissionAgent[user].Outbox[uuid.String()] = email
		} else {
			w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
		}
	} else {
		// For when we cannot make the email
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Read a particular email from the inbox or outbox
func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]
	uuid := vars["uuid"]

	// If the user in question does exist
	if account, ok := mailSubmissionAgent[user]; ok {
		// Has the type of box been specified?
		if box == "inbox" || box == "outbox" {

			var chosenbox map[string]Email // Removes the necesity to duplicate fucntions for both the inbox and outbox

			if box == "inbox" {
				chosenbox = account.Inbox
			} else {
				chosenbox = account.Outbox
			}
			// If the email in question does exist //PROBLEM: Replace line below if everything else works
			if _, ok := chosenbox[uuid]; ok {
				email := chosenbox[uuid]
				// If there are no errors converting it to JSON
				if enc, err := json.Marshal(email); err == nil {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(enc)) // Write the output onto the response
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				// Cannot find the email
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			// Must state whether targeting the inbox or outbox
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// Cannot find the user
		w.WriteHeader(http.StatusNotFound)
	}
}

// Delete a particular email from the outbox
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]
	uuid := vars["uuid"]

	// If the user in question does exist
	if account, ok := mailSubmissionAgent[user]; ok {

		// Has the type of box been specified?
		if box == "inbox" || box == "outbox" {

			var chosenbox map[string]Email // Removes the necesity to duplicate fucntions for both the inbox and outbox
			if box == "inbox" {
				chosenbox = account.Inbox
			} else {
				chosenbox = account.Outbox
			}
			// If the email in question does exist
			if _, ok := chosenbox[uuid]; ok {
				w.WriteHeader(http.StatusNoContent) // It all worked, but I have nothing else to say
				delete(chosenbox, uuid)
			} else {
				// Cannot find the email
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			// Must state whether targeting the inbox or outbox
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// Cannot find the user
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mailsubmissionagent/{user}", Create).Methods("POST")
	router.HandleFunc("/mailsubmissionagent/{user}/{box}", List).Methods("GET")
	router.HandleFunc("/mailsubmissionagent/{user}/{box}/{uuid}", Read).Methods("GET")
	router.HandleFunc("/mailsubmissionagent/{user}/{box}/{uuid}", Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func main() {
	mailSubmissionAgent = make(map[string]User)
	handleRequests()
}
