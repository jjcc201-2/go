package main

import (
	// Use JSON
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// An email is made out of these constituents
type Email struct {
	From    string
	To      string
	Message string
}

// var emails []Email
// var emailKeys []string

/*
func DeleteRequest(port int64, reqType string) {
	local := "http://localhost:"
	baseURL := local + strconv.FormatInt(port, 10) + "/mailsubmissionagent/Alice/outbox"
	client := &http.Client{}

}
*/
/*
func MakeRequest(port int64, reqType string, emailUuid string, reqName string) (*Response, error) {
	local := "http://localhost:"
	baseURL := local + strconv.FormatInt(port, 10) + "/mailsubmissionagent/Alice/outbox" + emailUuid
	client := &http.Client{}

	if req, createRequestErr := http.NewRequest(reqType, baseURL, nil); createRequestErr == nil {

		if resp, sendRequestErr := client.Do(req); sendRequestErr == nil {
			return resp, nil
		} else {
			return nil, (reqName + " failed with" + sendRequestErr.Error() + "\n")
		}
	} else {
		return nil, errors.New(reqName + " failed with" + createRequestErr.Error() + "\n")
	}

}
*/

func DeleteRequest(emailUuid string) error {
	url := "http://localhost:4000/MSA/Alice/outbox/" + emailUuid
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
	url := "http://localhost:4000/MSA/Alice/outbox/" + emailUuid
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
	url := "http://localhost:4000/MSA/Alice/outbox"
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
Make a post request, moving the
*/
func PostRequest(email Email) error {
	url := "http://localhost:4000/MSA/Alice/inbox"
	client := &http.Client{}
	// Can we convert it to JSON?
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

/*
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/mailtransferagent/{user}/{box}", Create).Methods("POST")
	router.HandleFunc("/MTA/{user}/{box}", List).Methods("GET")
	router.HandleFunc("/MTA/{user}/{box}/{uuid}", Read).Methods("GET")
	//router.HandleFunc("/mailsubmissionagent/{user}/{box}/{uuid}", Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4001", router))
}
*/
func main() {
	batchMove()

}

func batchMove() {
	for true {
		// var abc Email
		//emailsToTransfer := []string{}
		if emailList, listErr := ListRequest(); listErr == nil {

			for i := 0; i < len(emailList); i++ {

				if email, readErr := ReadRequest(emailList[i]); readErr == nil {

					if postError := PostRequest(email); postError == nil {

						if deleteErr := DeleteRequest(emailList[i]); deleteErr == nil {
							// Success
						} else {
							fmt.Println("Delete failed due to: " + deleteErr.Error())
						}
					} else {
						fmt.Println("Post failed due to " + postError.Error())
					}
				} else {
					fmt.Println("Read failed due to " + readErr.Error())
				}
			}
		} else {
			fmt.Println("List failed due to " + listErr.Error())
		}
		time.Sleep(7 * time.Second)
	}
}

/*
func sendToInbox() {
	emailsToTransfer := []string{}
	if emailList, listErr := ListRequest(4000, "GET"); listErr == nil {
		for i := 0; i < len(emailList); i++ {

			if email, readErr := ReadRequest(4000, "GET", emailList[i]); readErr == nil {
				emailsToTransfer = append(emailsToTransfer, email)

				if deleteErr := DeleteRequest(4000, "DELETE", emailList[i]); deleteErr == nil {
					fmt.Println("Delte successful")
				} else {
					fmt.Println("Delete failed due to: " + deleteErr.Error())
				}
				fmt.Println(email)
			} else {
				fmt.Println("Read failed due to: " + readErr.Error())
			}
		}
	} else {
		fmt.Println("List failed due to: " + listErr.Error())
	}
}
*/
/*
func main() {
	for true {
		var emailList []string
		fmt.Println("Before")
		fmt.Println(emailList)
		emailList = append(emailList, "hello")
		fmt.Println("After")
		fmt.Println(emailList)
		time.Sleep(1 * time.Second)
	}
}
*/
