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

func appendToFile(filePath string, msg string) {
	f, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	f.WriteString(msg + "\n")
	f.Close()
}

func readFileContent(filePath string) string {
	content, _ := ioutil.ReadFile(filePath)

	return string(content)
}