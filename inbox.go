package main 
import (
	"encoding/json" 			// Use JSON
	"github.com/gorilla/mux"    // A request router and dispatcher for matching incoming requests against a list of registered routes
	"log" 						// Simple logger 
	"net/http" 					// Provides HTTP client and server implementations. GET, POST, HEAD and PostForm
	"github.com/google/uuid"
)

/*
1. READ list of emails in inbox for a user http://localhost:8888/{user}
2. READ specific email					   http://localhost:8888/{user}/{emailKey}
3. DELETE specific email 				   http://localhost:8888/{user}/{emailKey}
4. Add email (from another user's OB)      http://localhost:8888/{user}
*/


type Email struct { // Type is used to refer to the struct afterwards
	from     string
	to 	     string
	message  string
}


type Inbox struct {
	email 	Email
}

var emails map[string] Email // Used as a wimpy way of "storing" replies on the server
var inbox map[string] Inbox

// List emails from a user's inbox
func List( w http.ResponseWriter, r *http.Request ) {

}

// Read specifc email from a user's inbox
func Read( w http.ResponseWriter, r *http.Request ) {

}

// Delete specific email from a user's inbox
func Delete( w http.ResponseWriter, r *http.Request ) {

}



// Add email that was sent to a user
func Add( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	user := vars[ "user" ]
	var email Email

	// DO WE NEED TO CHECK IF USER EXISTS??

	if uuid, err := uuid.NewUUID(); err == nil { // If successful generating a new UUID for email
		w.Header().Set("Location", r.Host+"/inbox/user/"+uuid.String())
		
	
	
	
	} else {
		w.WriteHeader(http.StatusInternalServerError) // Failed to add the email 
	}

	
}

vars := mux.Vars( r ) // The variables from the request
	user := vars[ "user" ] // The variables from the request
	decoder := json.NewDecoder( r.Body ) // Decode the body, which is stored as a JSON structure
	
	var email Email
	
	if err := decoder.Decode( &reply ); err == nil { // If there are no errors
		w.WriteHeader( http.StatusCreated )
		replies[ user ] = reply // Store away the reply. A very weak way to do so, but persistence not necessary for the CA?
	} else {
		w.WriteHeader( http.StatusBadRequest ) // If there is an error with the JSON, send back a bad status request
	}



func handleRequests() {
	router := mux.NewRouter().StrictSlash( true ) // 
	router.HandleFunc( "/inbox/{user}", List ).Methods( "GET" )
	router.HandleFunc( "/inbox/{user}/{emailKey}", Read ).Methods( "GET" )
	router.HandleFunc( "/inbox/{user}", Add ).Methods( "POST" )
	router.HandleFunc( "/inbox/{user}/{emailKey}", Delete ).Methods( "DELETE" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}


func main() {
	emails = make( map[ string ] Emails)
	handleRequests()
}


















func Create( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r ) // The variables from the request
	user := vars[ "user" ] // The variables from the request
	decoder := json.NewDecoder( r.Body ) // Decode the body, which is stored as a JSON structure
	
	var email Email
	
	if err := decoder.Decode( &reply ); err == nil { // If there are no errors
		w.WriteHeader( http.StatusCreated )
		replies[ user ] = reply // Store away the reply. A very weak way to do so, but persistence not necessary for the CA?
	} else {
		w.WriteHeader( http.StatusBadRequest ) // If there is an error with the JSON, send back a bad status request
	}
}



func Read( w http.ResponseWriter, r *http.Request ) {
    // The variables from the request
	vars := mux.Vars( r )
	
	// The variables from the request
	user := vars[ "user" ]
	
	// Look into the store of OOF auto replies to see if there is one for this particular user
	// When you look into these maps, you get two values back. A reply, and whether or not you can trust that reply
	if reply, ok := replies[ user ]; ok {
		w.WriteHeader( http.StatusOK )
		if enc, err := json.Marshal( reply ); err == nil { // If you have an error converting it to JSON
			w.Write( []byte( enc ) )		
		} else {
			w.WriteHeader( http.StatusInternalServerError )
		}
		
	// If reply isn't meaningful e.g. if Justin isn't in the reply list
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}


func Update( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	user := vars[ "user" ]
	
	if _, ok := replies[ user ]; ok { // Have we already got an OOO for this user?
		var reply Reply
		decoder := json.NewDecoder( r.Body )
		if err := decoder.Decode( &reply ); err == nil {
			w.WriteHeader( http.StatusCreated )
			replies[ user ] = reply
		} else {
			w.WriteHeader( http.StatusBadRequest ) // If we couldn't unMarshal the JSON
		}
	} else {
		w.WriteHeader( http.StatusNotFound ) // If we don't have an OOO for this user
	}
}


func Delete( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	user := vars[ "user" ]
	
	if _, ok := replies[ user ]; ok { // Have we already got an OOO for this user?
		w.WriteHeader( http.StatusNoContent ) // It all worked, but I have nothing else to say
		delete( replies, user )
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}


