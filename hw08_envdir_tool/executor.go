package main

import (
	"errors"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	for key, envValue := range env {
		if envValue.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				return 0
			}
		} else {
			if err := os.Setenv(key, envValue.Value); err != nil {
				return 0
			}
		}
	}

	command.Env = os.Environ()

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}
