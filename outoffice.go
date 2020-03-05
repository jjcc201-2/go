package main 
import (
	"encoding/json" 			// Use JSON
	"github.com/gorilla/mux"    // A request router and dispatcher for matching incoming requests against a list of registered routes
	"log" 						// Simple logger 
	"net/http" 					// Provides HTTP client and server implementations. GET, POST, HEAD and PostForm
)



type Reply struct { // Type is used to refer to the struct afterwards
	Summary string
}

var replies map[string] Reply // Used as a wimpy way of "storing" replies on the server






func Create( w http.ResponseWriter, r *http.Request ) {
	// The variables from the request
	vars := mux.Vars( r )
	
	// The variables from the request
	user := vars[ "user" ] 
	
	// Decode the body, which is stored as a JSON structure
	decoder := json.NewDecoder( r.Body )
	
	var reply Reply
	
	// If there are nil errors
	if err := decoder.Decode( &reply ); err == nil {
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



func handleRequests() {
	router := mux.NewRouter().StrictSlash( true ) // 
	router.HandleFunc( "/outoffice/{user}", Create ).Methods( "POST" )
	router.HandleFunc( "/outoffice/{user}", Read ).Methods( "GET" )
	router.HandleFunc( "/outoffice/{user}", Update ).Methods( "PUT" )
	router.HandleFunc( "/outoffice/{user}", Delete ).Methods( "DELETE" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}

func main() {
	replies = make( map[ string ] Reply)
	handleRequests()
}