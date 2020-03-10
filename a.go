package main

import (
	"fmt"
)

func main() {
	box := "otbox"
	if box == "inbox" || box == "outbox" {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
