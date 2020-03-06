package main

import (
	// Use JSON
	"encoding/json"
	"fmt"

	// Simple logger
	// Provides HTTP client and server implementations. GET, POST, HEAD and PostForm
	"github.com/google/uuid"
	// A request router and dispatcher for matching incoming requests against a list of registered routes
)

var rawJson = []byte(`[{"from":"me", "to":"you", "message":"hello"}]`)

type Email struct { // Type is used to refer to the struct afterwards
	from    string
	to      string
	message string
}

// Create email and add it to the user's outbox
func Create(rawJson []byte) {
	decoder := json.NewDecoder(rawJson)
	var email Email

	if uuid, err := uuid.NewUUID(); err == nil { // If successful generating a new UUID for email
		fmt.Println(uuid)
		// w.Header().Set("Location", r.Host+"/outbox/user/"+uuid.String())
		if err := decoder.Decode(&email); err == nil { //If no errors in decoding the message into the email object
			emails[user] = email
			fmt.Println("abc ", email)
		} else {
		}
	} else {
	}
}
