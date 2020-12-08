package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envVarName, envInfo := range env {
		if err := os.Unsetenv(envVarName); err != nil {
			fmt.Println(err)
			return -1
		}

		if !envInfo.NeedRemove {
			if err := os.Setenv(envVarName, envInfo.Value); err != nil {
				fmt.Println(err)
				return -1
			}
		}
	}

	name := cmd[0]
	args := cmd[1:]

	command := exec.Command(name, args...)
	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println(err)
	}

	return command.ProcessState.ExitCode()
}
