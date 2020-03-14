package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Use JSON

var bluebook map[string]string

func GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]

	// Does the server exist?
	if server, ok := bluebook[domain]; ok {
		// Can the result be marshalled?
		if enc, err := json.Marshal(server); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(enc))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // Server does not exist
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/bluebook/{domain}", GetAddress).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
	bluebook = make(map[string]string)
	bluebook["here.com"] = ":4001"
	bluebook["there.com"] = ":8001"
	handleRequests()
}
