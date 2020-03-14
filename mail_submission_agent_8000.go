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

// Each user has an inbox and an outbox, which are the storage containers for emails
type User struct {
	Inbox  map[string]Email
	Outbox map[string]Email
}

// Used to keep account of who is on the server
var userList []string

var mailSubmissionAgent map[string]User

/*
Helper function that returns the box specified with an API request
*/
func getBox(box string, account User) map[string]Email {
	if box == "inbox" {
		return account.Inbox
	} else if box == "outbox" {
		return account.Outbox
	} else {
		return nil
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Convert key array to JSON
	if enc, err := json.Marshal(userList); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(enc))
	} else {
		// JSON conversion error
		w.WriteHeader(http.StatusInternalServerError)
	}

}

/*
Function lists all of the IDs for emails stored in a user's inbox or a user's outbox
*/
func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]

	// Check if user exists
	if account, ok := mailSubmissionAgent[user]; ok {
		// Check if box is specified
		if chosenbox := getBox(box, account); chosenbox != nil {
			// Get all email IDs in the specified inbox or outbox
			keys := []string{}
			for emailUuid, _ := range chosenbox {
				keys = append(keys, emailUuid)
			}
			// Convert key array to JSON
			if enc, err := json.Marshal(keys); err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(enc))
			} else {
				// JSON conversion error
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			// Must state whether targeting the inbox or outbox
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// User does not exist
		w.WriteHeader(http.StatusNotFound)
	}
}

/*
Function creates an email based on the body of the request and the user and outbox specified
*/
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]
	decoder := json.NewDecoder(r.Body)
	var email Email

	// UUID to uniquely identify an email
	if uuid, err := uuid.NewUUID(); err == nil {
		// Check if user exists
		if mailSubmissionAgent[user].Inbox == nil || mailSubmissionAgent[user].Outbox == nil {
			// If not, create them and then store the new user in the list
			mailSubmissionAgent[user] = User{Inbox: make(map[string]Email), Outbox: make(map[string]Email)}
			userList = append(userList, user)
		}
		account := mailSubmissionAgent[user]
		// Check if box is specified
		if chosenbox := getBox(box, account); chosenbox != nil {
			// Try to decode the response into the email structure
			if err := decoder.Decode(&email); err == nil {
				w.WriteHeader(http.StatusCreated)
				chosenbox[uuid.String()] = email
			} else {
				// JSON conversion error
				w.WriteHeader(http.StatusBadRequest)
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

/*
Function reads a specific email from either the inbox or outbox by looking up the email ID specified
*/
func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]
	uuid := vars["uuid"]

	// Check if user exists
	if account, ok := mailSubmissionAgent[user]; ok {
		// Check if box is specified
		if chosenbox := getBox(box, account); chosenbox != nil {
			// If the email in question does exist
			if email, ok := chosenbox[uuid]; ok {
				// Convert email to JSON
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
		// User does not exist
		w.WriteHeader(http.StatusNotFound)
	}
}

/*
Function deletes a specific email from either the inbox or outbox by looking up the email ID specified
*/
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	box := vars["box"]
	uuid := vars["uuid"]

	// Check if user exists
	if account, ok := mailSubmissionAgent[user]; ok {
		// Check if box is specified
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

/*
Function handles all potential routes specified within the Mail Submission Agent
*/
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/MSA/{user}/{box}", Create).Methods("POST")
	router.HandleFunc("/MSA/{user}/{box}", List).Methods("GET")
	router.HandleFunc("/MSA/{user}/{box}/{uuid}", Read).Methods("GET")
	router.HandleFunc("/MSA/{user}/{box}/{uuid}", Delete).Methods("DELETE")
	router.HandleFunc("/MSA", GetUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	mailSubmissionAgent = make(map[string]User)
	handleRequests()
}
