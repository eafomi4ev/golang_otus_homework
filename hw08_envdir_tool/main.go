package main

import (
	"fmt"
	"os"
)

func main() {
	envVars, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Errorf("error occurred while reading envs from dir: %s; error is: %w", os.Args[1], err))
	}

	RunCmd(os.Args[2:], envVars)
}
