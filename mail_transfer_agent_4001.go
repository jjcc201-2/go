package main

import (
	// Use JSON
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// An email is made out of these constituents
type Email struct {
	From    string
	To      string
	Message string
}

/*
Helper function that builds an HTTP request for the MSA based on the parameters
inputted.
*/
func RequestBuilder(reqType string, url string, enc []byte) (*http.Request, error) {
	if reqType == "POST" {
		req, createRequestErr := http.NewRequest(reqType, url, bytes.NewBuffer(enc))
		return req, createRequestErr // If an error occurs, it will propogate back
	}
	if reqType == "GET" || reqType == "DELETE" {
		req, createRequestErr := http.NewRequest(reqType, url, nil)
		return req, createRequestErr
	}
	return nil, errors.New("A request type must be specified")
}

/*
Helper function that carries out the fundamental stages for every request, no matter
the type. It will build the HTTP request, send it off and check for the resulting
status code returned by the MSA.
*/
func GeneralRequest(reqType string, url string, enc []byte) (*http.Response, error) {
	client := &http.Client{}
	// Check if request can be built
	if req, createRequestErr := RequestBuilder(reqType, url, enc); createRequestErr == nil {
		// Check if there are no problems sending the request
		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {
			// Check if it's a positive result
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				return resp, nil
			} else {
				return nil, errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(http.StatusText(resp.StatusCode)))
			}
		} else {
			return nil, errors.New(sendRequestErr.Error())
		}
	} else {
		return nil, errors.New(createRequestErr.Error())
	}
}

/*
Function makes a DELETE request to the MSA in order to delete a specific
email in a user's outbox
*/
func DeleteRequest(emailUuid string, user string) error {
	url := "http://localhost:4000/MSA/" + user + "/outbox/" + emailUuid
	// Make POST request with specified URL
	if _, requestErr := GeneralRequest("DELETE", url, nil); requestErr == nil {
		return nil // Delete completed, return nil errors
	} else {
		return requestErr
	}
}

/*
Function makes a GET request to the MSA in order to obtain an email from a
specific user's outbox.
*/
func ReadRequest(emailUuid string, user string) (Email, error) {
	url := "http://localhost:4000/MSA/" + user + "/outbox/" + emailUuid
	var email Email
	// Make GET request with specified URL
	if resp, requestErr := GeneralRequest("GET", url, nil); requestErr == nil {
		decoder := json.NewDecoder(resp.Body)

		if err := decoder.Decode(&email); err == nil {
			return email, nil
		} else {
			return email, errors.New(err.Error())
		}
	} else {
		return email, requestErr
	}
}

/*
Function makes a GET request to the MSA in order to obtain email IDs of
all emails situated in a specific users outbox
*/
func ListRequest(user string) ([]string, error) {
	url := "http://localhost:4000/MSA/" + user + "/outbox/"
	var emailKeys []string
	// Make GET request with specified URL
	if resp, requestErr := GeneralRequest("GET", url, nil); requestErr == nil {
		decoder := json.NewDecoder(resp.Body)

		if decodeErr := decoder.Decode(&emailKeys); decodeErr == nil {
			return emailKeys, nil

		} else {
			return nil, errors.New(decodeErr.Error())
		}
	} else {
		return nil, requestErr
	}
}

/*
Send an email to another mail transfer agent on another servers
*/
func SendToMTAServer(email Email, address string) error {
	url := "http://localhost" + address + "/MTA"
	// Convert email to JSON
	if enc, jsonConversionErr := json.Marshal(email); jsonConversionErr == nil {
		// Make POST request with specified URL
		if _, requestErr := GeneralRequest("POST", url, enc); requestErr == nil {
			return nil // Request complete. Nothing else to do aside from return no errors
		} else {
			return requestErr
		}
	} else {
		return jsonConversionErr
	}
}

/*
Function makes a GET request to the MSA in order to obtain all users currently part of the email server
*/
func ObtainAllUsers() ([]string, error) {
	url := "http://localhost:4000/MSA"
	var allUsers []string // How users ids will be stored

	if resp, requestErr := GeneralRequest("GET", url, nil); requestErr == nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&allUsers); err == nil {
			return allUsers, nil
		} else {
			return nil, errors.New(err.Error())
		}
	} else {
		return nil, requestErr
	}
}

/*
Function makes a GET request to the MSA in order to locate the
server which an email is destined for
*/
func ObtainBluebookAddress(domain string) (string, error) {
	url := "http://localhost:3000/bluebook/" + domain
	var address string
	if resp, requestErr := GeneralRequest("GET", url, nil); requestErr == nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&address); err == nil {
			return address, nil
		} else {
			return "", errors.New(err.Error())
		}
	} else {
		return "", requestErr
	}
}

/*
Function is triggered by an MTA from another server. It gives this MTA
an email, which is then posted into the the destined user's outbox
*/
func PostToMSA(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var email Email
	//If no errors in decoding the message into the email object
	if err := decoder.Decode(&email); err == nil {
		receivingUser := email.To
		url := "http://localhost:4000/MSA/" + receivingUser + "/inbox"
		// Convert back to JSON
		if enc, jsonConversionErr := json.Marshal(email); jsonConversionErr == nil {
			// Make POST request to request helper function
			if _, requestErr := GeneralRequest("POST", url, enc); requestErr == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

/*
Function will perform multiple requests to the MSA, which allows it to
get access all users on the server and obtain their emails. It will then
send them off to the servers that contain the recipients of the emails.
This action is performed periodically.
*/
func batchMove() {
	for {
		// Get all users on the server
		if userList, userErr := ObtainAllUsers(); userErr == nil {
			for i := 0; i < len(userList); i++ {
				user := userList[i]
				// Check if they have any emails in their outbox
				if emailList, listErr := ListRequest(user); listErr == nil {
					for j := 0; j < len(emailList); j++ { // Goes through each email individually
						if email, readErr := ReadRequest(emailList[j], user); readErr == nil {
							split := strings.Split(email.To, "@") // This gets the specific domain name
							_, domain := split[0], split[1]
							// Locate the address of the user destined for the email
							if address, addressError := ObtainBluebookAddress(domain); addressError == nil {
								// Send it the email to the recipient's server and then delete the original from the outbox
								if sendError := SendToMTAServer(email, address); sendError == nil {
									if deleteErr := DeleteRequest(emailList[j], user); deleteErr == nil {
										// Success. No more action needed
									} else {
										fmt.Println("Delete failed due to: " + deleteErr.Error())
									}
								} else {
									fmt.Println("Post failed due to " + sendError.Error())
								}
							} else {
								fmt.Println("Could not obtain address due to " + addressError.Error())
							}
						} else {
							fmt.Println("Read failed due to " + readErr.Error())
						}
					}
				} else {
					fmt.Println("Email listing failed due to " + listErr.Error())
				}
			}
		} else {
			fmt.Println("User fetch failed due to " + userErr.Error())
		}
		time.Sleep(5 * time.Second) // Sleep for 5 seconds before next iteration
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/MTA", PostToMSA).Methods("POST")
	log.Fatal(http.ListenAndServe(":4001", router))
}

func main() {
	go batchMove() // Run this on another thread to allow program to also listen for requests
	handleRequests()
}
