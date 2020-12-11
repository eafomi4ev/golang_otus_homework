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

	if err := scanner.Err(); err != nil {
		return "", errors.Wrap(err, "")
	}
	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "")
	}

	return firstLine, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if strings.Contains(dir, "=") {
		return nil, errors.Wrap(ErrIncorrectFileName, "")
	}

	dirContent, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var dirFilesInfo []os.FileInfo
	for _, item := range dirContent {
		if !item.IsDir() {
			dirFilesInfo = append(dirFilesInfo, item)
		}
	}

	envMap := make(Environment)

	for _, dirFileInfo := range dirFilesInfo {
		if dirFileInfo.Size() == 0 {
			envMap[dirFileInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		firstLine, err := readFirstLineFromFile(filepath.Join(dir, dirFileInfo.Name()))
		if err != nil {
			return nil, err
		}

		firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")
		firstLine = strings.TrimRight(firstLine, " \t")

		envMap[dirFileInfo.Name()] = EnvValue{Value: firstLine, NeedRemove: false}
	}

	return envMap, nil
}
