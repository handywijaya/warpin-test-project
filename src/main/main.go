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
	http.HandleFunc("/flush", flush)
	http.HandleFunc("/try", try)

	http.ListenAndServe(":8090", nil)
}

// This API will save every messages sent with key "msg"
// and store them into a file
func send(w http.ResponseWriter, req *http.Request) {
	msg := req.URL.Query()["msg"][0]

	// use mutex to prevent race condition
	mutex.Lock()
	err := appendToFile(getFilePath(dataFileName), msg)
	mutex.Unlock()

	// error checking after operation
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return response with error checking
	_, err = fmt.Fprintf(w, "Data received!\nUse /get API to retrieve all the messages sent.")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will read file where the message sent is stored
// and then print the content
func get(w http.ResponseWriter, req *http.Request) {
	// use mutex to prevent race condition
	mutex.Lock()
	content := readFileContent(getFilePath(dataFileName))
	mutex.Unlock()

	// return response with error checking
	_, err := fmt.Fprintf(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will delete all messages sent from the data file
func flush(w http.ResponseWriter, req *http.Request) {
	// use mutex to prevent race condition
	mutex.Lock()
	err := flushFileContent(getFilePath(dataFileName))
	mutex.Unlock()

	// error checking after operation
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return response with error checking
	_, err = fmt.Fprintf(w, "Data cleaned! Use /send API to send message.")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will send the HTML page for client
// to connect through websocket
func try(w http.ResponseWriter, req *http.Request) {
	content := readFileContent(getFilePath(htmlFileName))

	// return response with error checking
	_, err := fmt.Fprintf(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}