package main

import "fmt"

type Email struct { // Type is used to refer to the struct afterwards
	emailKey string
	from     string
	to       string
	message  string
}

func main() {

	emails := make(map[string]Email)

	emails["a"] = Email{emailKey: "1", from: "a@b.com", to: "z@y.com", message: "hello"}
	emails["b"] = Email{emailKey: "2", from: "c@d.com", to: "x@w.com", message: "what is up"}
	emails["c"] = Email{emailKey: "3", from: "e@f.com", to: "t@u.com", message: "goodbye"}

	for k, v := range emails {
		fmt.Println("Key: ", k)
		fmt.Println("from: ", v.from)
		fmt.Println("to: ", v.to)
		fmt.Println("message: ", v.message)
		fmt.Println("\n")
	}

	/*
		grades := make(map[string]float32) // A string key and a float32 value. Make is used to actually add values

		grades["Timmy"] = 42
		grades["Jess"] = 92
		grades["Sam"] = 67

		fmt.Println(grades) // Print out entire map

		TimsGrade := grades["Timmy"]
		fmt.Println(TimsGrade)

		delete(grades, "Timmy")
		fmt.Println(grades) // Print out entire map

		for k, v := range grades {
			fmt.Println(k, ":", v)
		}
	*/

}
