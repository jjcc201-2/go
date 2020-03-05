package main
import (
	"fmt"
	"net/http"
)

func main() {
	url    := "http://localhost:8888/outoffice/David"
	client := &http.Client {}
	if req, err := http.NewRequest( "DELETE", url, nil); err == nil {
		if _, err1 := client.Do( req ); err1 == nil {
			// Nothing
		} else {
			fmt.Printf( "DELETE failed with %s\n", err1 )
		}
	} else {
		fmt.Printf( "Delete failed with %s\n", err )
	}
}