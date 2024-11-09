package main

import (
	"fmt"
	"os"
)

var (
	envDir, shell, scriptPath string
	scriptArgs                []string
)

func init() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go-envdir <envdir_path> <shell> <script_path> [script_args...]")
		os.Exit(2)
	}

	envDir = os.Args[1]
	shell = os.Args[2]
	scriptPath = os.Args[3]
	scriptArgs = os.Args[4:]
}

func main() {
	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := append([]string{shell, scriptPath}, scriptArgs...)
	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
