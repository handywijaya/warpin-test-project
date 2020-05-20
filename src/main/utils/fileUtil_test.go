package utils

import (
	"os"
	"testing"
)

func assert(t *testing.T, result interface{}, expected interface{}) {
	if result != expected {
		t.Fail()
	}
}

func TestGetBaseDirectory(t *testing.T) {
	result := GetBaseDirectory()
	expected,_ := os.Getwd()

	assert(t, result, expected)
}

func TestGetFilePath(t *testing.T) {
	fileName := "test.txt"

	result := GetFilePath(fileName)

	dirPath, _ := os.Getwd()
	expected := dirPath + "/resources/" + fileName

	assert(t, result, expected)
}
