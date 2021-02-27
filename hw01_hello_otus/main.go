package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	A := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(A)
}

// func Returner() string {
// 	return A
// }
