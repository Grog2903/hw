package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		value, needRemove, err := processEnvFile(dir, file.Name())
		if err != nil {
			return nil, err
		}

		env[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: needRemove,
		}
	}

	return env, nil
}

func processEnvFile(dir string, fileName string) (value string, needRemove bool, err error) {
	filePath := dir + "/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", false, err
	}
	if info.Size() == 0 {
		return "", true, nil
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, " \t")
		line = strings.ReplaceAll(line, "\x00", "\n")

		return line, false, nil
	}

	if err := scanner.Err(); err != nil {
		return "", false, err
	}

	return "", false, nil
}
