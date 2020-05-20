package utils

import (
	"io/ioutil"
	"os"
)

func GetBaseDirectory() string {
	path, _ := os.Getwd()
	return path
}

func GetFilePath(fileName string) string {
	return GetBaseDirectory() + "/" + "resources" + "/" + fileName
}

func AppendToFile(filePath string, msg string) error {
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

func ReadFileContent(filePath string) string {
	content, _ := ioutil.ReadFile(filePath)

	return string(content)
}

func FlushFileContent(filePath string) error {
	_, err := os.Create(filePath)
	return err
}