package main

import (
	"os"
)

func main() {
	envVars, err := ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}

	RunCmd(os.Args[2:], envVars)
}
