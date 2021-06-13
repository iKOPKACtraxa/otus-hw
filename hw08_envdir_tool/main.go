package main

import (
	"fmt"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("reading of environment from dir has got an error:", err)
	}
	returnCode, err := RunCmd(os.Args, env)
	if err != nil {
		fmt.Println("running of cmd has got an error:", err)
	}
	os.Exit(returnCode)
}
