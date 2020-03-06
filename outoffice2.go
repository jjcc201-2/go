package main

import (
	"fmt" // Provides HTTP client and server implementations. GET, POST, HEAD and PostForm
)

type Reply struct { // Type is used to refer to the struct afterwards
	Summary string
}

type Replies map[string]Reply

var outbox map[string]Replies

/*
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	decoder := json.NewDecoder(r.Body)

	var reply Reply

	if err := decoder.Decode(&reply); err == nil {
		w.WriteHeader(http.StatusCreated)
		outbox[user]["1"] = reply
	} else {
		w.WriteHeader(http.StatusBadRequest) // If there is an error with the JSON, send back a bad status request
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true) //
	router.HandleFunc("/outoffice/{user}", Create).Methods("POST")
	log.Fatal(http.ListenAndServe(":8888", router))
}
*/

func main() {

	// outbox := make(map[string]map[string]Reply)
	outbox := make(map[string]Replies)

	if outbox["Bob"] == nil {
		fmt.Println("nil")
		outbox["Bob"] = make(map[string]Reply)
	}
	outbox["Bob"]["1"] = Reply{"Help Me"}
	outbox["Bob"]["2"] = Reply{"Fixed"}
	outbox["Bob"]["3"] = Reply{"Wham bam koblam"}
	fmt.Println(outbox["Bob"]["1"])
	fmt.Println(outbox["Bob"]["2"])
	fmt.Println(outbox["Bob"]["3"])
	// handleRequests()
}
