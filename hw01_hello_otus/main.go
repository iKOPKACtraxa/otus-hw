package main

import (
	"fmt"
	"os"

	"golang.org/x/example/stringutil"
)

func main() {
	fmt.Println(stringutil.Reverse("Hello, OTUS!"))
	gg := os.Stdout
	fmt.Println(gg)

}
