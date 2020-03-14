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

var userMTA string

var sentEmails []Email

func DeleteRequest(emailUuid string) error {

	url := "http://localhost:4000/MSA/" + userMTA + "/outbox/" + emailUuid
	client := &http.Client{}

	if req, createRequestErr := http.NewRequest("DELETE", url, nil); createRequestErr == nil {

		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				return nil // Delete completed, send back no errors
			} else {
				return errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(http.StatusText(resp.StatusCode)))
			}
		} else {
			return errors.New(sendRequestErr.Error())
		}
	} else {
		return errors.New(createRequestErr.Error())
	}

}

/*

 */
func ReadRequest(emailUuid string) (Email, error) {
	url := "http://localhost:4000/MSA/" + userMTA + "/outbox/" + emailUuid
	client := &http.Client{}
	var email Email
	if req, createRequestErr := http.NewRequest("GET", url, nil); createRequestErr == nil {

		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

				decoder := json.NewDecoder(resp.Body)

				if err := decoder.Decode(&email); err == nil {
					return email, nil
				} else {
					return email, errors.New(err.Error())
				}
			} else {
				return email, errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(http.StatusText(resp.StatusCode)))
			}
		} else {
			return email, errors.New(sendRequestErr.Error())
		}
	} else {
		return email, errors.New(createRequestErr.Error())
	}
}

/*
Make a List request to the Mail Submission Agent, storing each email ID for reference
*/
func ListRequest() ([]string, error) {
	url := "http://localhost:4000/MSA/" + userMTA + "/outbox/"
	client := &http.Client{}
	var emailKeys []string

	if req, createRequestErr := http.NewRequest("GET", url, nil); createRequestErr == nil {

		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

				decoder := json.NewDecoder(resp.Body)

				if err := decoder.Decode(&emailKeys); err == nil {
					return emailKeys, nil

				} else {
					return nil, errors.New(err.Error())
				}
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
Send an email to another mail transfer agent on another servers
*/
func SendToMTAServer(email Email, address string) error {
	url := "http://localhost" + address + "/MTA"
	client := &http.Client{}

	if enc, jsonConversionErr := json.Marshal(email); jsonConversionErr == nil {

		if req, createRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(enc)); createRequestErr == nil {

			if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

				if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

					return nil

				} else {
					return errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(http.StatusText(resp.StatusCode)))
				}
			} else {
				return errors.New(sendRequestErr.Error())
			}
		} else {
			return errors.New(createRequestErr.Error())
		}
	} else {
		return errors.New(jsonConversionErr.Error())
	}
}

func ObtainBluebookAddress(domain string) (string, error) {
	url := "http://localhost:3000/bluebook/" + domain
	client := &http.Client{}
	var address string
	if req, createRequestErr := http.NewRequest("GET", url, nil); createRequestErr == nil {

		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

				decoder := json.NewDecoder(resp.Body)

				if err := decoder.Decode(&address); err == nil {
					fmt.Println("Here We Go")
					fmt.Println(address)
					return address, nil
				} else {
					return "", errors.New(err.Error())
				}
			} else {
				return "", errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(http.StatusText(resp.StatusCode)))
			}
		} else {
			return "", errors.New(sendRequestErr.Error())
		}
	} else {
		return "", errors.New(createRequestErr.Error())
	}

}

/*

 */
func PostToMSA(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var email Email

	//If no errors in decoding the message into the email object
	if err := decoder.Decode(&email); err == nil {

		receivingUser := email.To
		url := "http://localhost:4000/MSA/" + receivingUser + "/inbox"
		client := &http.Client{}

		// Convert back to JSON
		if enc, jsonConversionErr := json.Marshal(email); jsonConversionErr == nil {

			if req, createRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(enc)); createRequestErr == nil {

				if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {

					if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

						w.WriteHeader(http.StatusOK)

					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func batchMove() {
	for {
		if emailList, listErr := ListRequest(); listErr == nil {

			for i := 0; i < len(emailList); i++ {

				if email, readErr := ReadRequest(emailList[i]); readErr == nil {

					// Get the domain name
					split := strings.Split(email.To, "@")
					_, domain := split[0], split[1]

					if address, addressError := ObtainBluebookAddress(domain); addressError == nil {

						if sendError := SendToMTAServer(email, address); sendError == nil {

							if deleteErr := DeleteRequest(emailList[i]); deleteErr == nil {
								// Success
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
			// Nothing. User might not have been instanciated yet
		}
		time.Sleep(7 * time.Second)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/MTA", PostToMSA).Methods("POST")
	log.Fatal(http.ListenAndServe(":4001", router))
}

func main() {
	userMTA = "alice@here.com"
	go batchMove()
	handleRequests()

}
