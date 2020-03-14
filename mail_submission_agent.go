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

// Helper function to return the correct box specified with the API request
func getBox(box string, account User) map[string]Email {
	if box == "inbox" {
		return account.Inbox
	} else if box == "outbox" {
		return account.Outbox
	} else {
		return nil
	}
}

// List emails in a user's inbox or outbox
func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]

	if account, ok := mailSubmissionAgent[user]; ok { // Does the user exist?

		if chosenbox := getBox(box, account); chosenbox != nil { // Has the box been specified?

			keys := []string{}
			for emailUuid, _ := range chosenbox {
				keys = append(keys, emailUuid)
			}

			if enc, err := json.Marshal(keys); err == nil { // If you have an error converting it to JSON
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
	box := vars["box"]
	decoder := json.NewDecoder(r.Body)
	var email Email

	// UUID to uniquely identify an email
	if uuid, err := uuid.NewUUID(); err == nil {

		// Does the user exist?
		if mailSubmissionAgent[user].Inbox == nil || mailSubmissionAgent[user].Outbox == nil {
			// If not, create them
			mailSubmissionAgent[user] = User{Inbox: make(map[string]Email), Outbox: make(map[string]Email)}
		}
		account := mailSubmissionAgent[user]

		if chosenbox := getBox(box, account); chosenbox != nil { // Has the box been specified?

			if err := decoder.Decode(&email); err == nil { //If no errors in decoding the message into the email object

				w.WriteHeader(http.StatusCreated)
				chosenbox[uuid.String()] = email
			} else {
				w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
			}
		} else {
			// Must state whether targeting the inbox or outbox
			w.WriteHeader(http.StatusBadRequest)
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
		if chosenbox := getBox(box, account); chosenbox != nil {
			// If the email in question does exist
			if email, ok := chosenbox[uuid]; ok {
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
		if chosenbox := getBox(box, account); chosenbox != nil { // Has the box been specified?
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
	router.HandleFunc("/MSA/{user}/{box}", Create).Methods("POST")
	router.HandleFunc("/MSA/{user}/{box}", List).Methods("GET")
	router.HandleFunc("/MSA/{user}/{box}/{uuid}", Read).Methods("GET")
	router.HandleFunc("/MSA/{user}/{box}/{uuid}", Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func main() {
	mailSubmissionAgent = make(map[string]User)
	handleRequests()
}
