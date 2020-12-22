package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]EnvValue

var ErrIncorrectFileName = errors.New("file name contains \"=\"")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readFirstLineFromFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()

	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "")
	}
	if err := scanner.Err(); err != nil {
		return "", errors.Wrap(err, "")
	}

	return firstLine, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirContent, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	envMap := make(Environment)
	for _, envFile := range dirContent {
		if envFile.IsDir() {
			continue
		}

		if strings.Contains(envFile.Name(), "=") {
			return nil, errors.Wrap(ErrIncorrectFileName, "")
		}

		if envFile.Size() == 0 {
			envMap[envFile.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		firstLine, err := readFirstLineFromFile(filepath.Join(dir, envFile.Name()))
		if err != nil {
			return nil, err
		}

		firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")
		firstLine = strings.TrimRight(firstLine, " \t")

		envMap[envFile.Name()] = EnvValue{Value: firstLine, NeedRemove: false}
	}

	return envMap, nil
}
