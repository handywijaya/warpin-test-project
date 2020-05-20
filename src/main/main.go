package main

import (
	"fmt"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}

func main() {
	http.HandleFunc("/send", send)
	http.HandleFunc("/get", get)
	http.HandleFunc("/try", try)

	http.ListenAndServe(":8090", nil)
}

// This API will save every messages sent with key "msg"
// and store them into a file
func send(w http.ResponseWriter, req *http.Request) {
	// use mutex
	msg := req.URL.Query()["msg"][0]

	mutex.Lock()
	appendToFile(getFilePath(dataFileName), msg)
	mutex.Unlock()

	fmt.Fprintf(w, "OK\n")
}

// This API will read file where the message sent is stored
// and then print the content
func get(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	content := readFileContent(getFilePath(dataFileName))
	mutex.Unlock()

	fmt.Fprintf(w, content)
}

// This API will send the HTML page for client
// to connect through websocket
func try(w http.ResponseWriter, req *http.Request) {
	content := readFileContent(getFilePath(htmlFileName))

	fmt.Fprintf(w, content)
}