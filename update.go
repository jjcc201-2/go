package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
Basic communication between services. Function is making a POST request to another Microservice
*/
func main() {
	url := "http://localhost:8888/outoffice/David"
	client := &http.Client {}

	data := map[ string ] string{ "Summary" : "Updated reply through update.go"}

	
	if enc, err:= json.Marshal(data); err == nil {
		// Can we build a PUT request?
		if req, err1 := http.NewRequest( "PUT", url, bytes.NewBuffer( enc ) ); err1 == nil {
			// Can we give to a HTTP client?
			if resp, err2 := client.Do( req ); err2 == nil {
				// Would it give us a result that we'd be interested in?
				if body, err3 := ioutil.ReadAll( resp.Body ); err3 == nil {
					// Print out reponse
					fmt.Println( string( body ) )
				} else {
					fmt.Printf( "POST failed with %s\n", err3 )
				}
			} else {
				fmt.Printf( "POST failed with %s\n", err2 )
			}
		} else {
			fmt.Printf( "POST failed with %s\n", err1 )
		}
	} else {
		fmt.Printf( "POST failed with %s\n", err )
	}
}