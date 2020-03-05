package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url    := "http://localhost:8888/outoffice/David"
	client := &http.Client {}
	// Can we build a GET request?
	if req, err :=  http.NewRequest( "GET", url, nil ); err == nil {
		// Can we give to a HTTP client?
		if resp, err1 := client.Do( req ); err1== nil {
			// Would it give us a result that we'd be interested in?
			if body, err2 := ioutil.ReadAll( resp.Body ); err2 == nil {
				// Print out reponse
				fmt.Println( string( body ) )
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