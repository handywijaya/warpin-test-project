package main

import (
	"io/ioutil"
	"os"
)

func getBaseDirectory() string {
	path, _ := os.Getwd()
	return path
}

func getFilePath(fileName string) string {
	return getBaseDirectory() + "/resources/" + fileName
}

func appendToFile(filePath string, msg string) error {
	f, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	_, err := f.WriteString(msg + "\n")
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	// return no error
	return nil
}

func readFileContent(filePath string) string {
	content, _ := ioutil.ReadFile(filePath)

	return string(content)
}

func flushFileContent(filePath string) error {
	_, err := os.Create(getFilePath(dataFileName))
	return err
}