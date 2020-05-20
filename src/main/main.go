package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
	"main/const"
	"main/pubsub"
	"main/utils"
	"net/http"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var broker = pubsub.Broker{Subscribers: make(map[string]chan string)}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/home", homeAPI)
	http.HandleFunc("/send", sendAPI)
	http.HandleFunc("/get", getAPI)
	http.HandleFunc("/flush", flushAPI)
	http.HandleFunc("/ws", wsPageAPI)
	http.HandleFunc("/ws-connection", wsConnectionAPI)

	http.ListenAndServe(":8090", nil)
}

// This API will send an instruction how to use this software
func homeAPI(w http.ResponseWriter, req *http.Request) {
	content := "Please read README.md."

	// return response with error checking
	_, err := fmt.Fprintf(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will save every messages sent with key "msg"
// and store them into a file
func sendAPI(w http.ResponseWriter, req *http.Request) {
	msg := req.URL.Query()["msg"][0]

	// use mutex to prevent race condition
	mutex.Lock()
	err := utils.AppendToFile(utils.GetFilePath(_const.DataFileName), msg)
	mutex.Unlock()

	broker.PublishMessage(msg)

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
func getAPI(w http.ResponseWriter, req *http.Request) {
	// use mutex to prevent race condition
	mutex.Lock()
	content := utils.ReadFileContent(utils.GetFilePath(_const.DataFileName))
	mutex.Unlock()

	// return response with error checking
	_, err := fmt.Fprintf(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will delete all messages sent from the data file
func flushAPI(w http.ResponseWriter, req *http.Request) {
	// use mutex to prevent race condition
	mutex.Lock()
	err := utils.FlushFileContent(utils.GetFilePath(_const.DataFileName))
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
func wsPageAPI(w http.ResponseWriter, req *http.Request) {
	content := utils.ReadFileContent(utils.GetFilePath(_const.HtmlFileName))

	// return response with error checking
	_, err := fmt.Fprintf(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This API will have a websocket connection to client
// to retrieve message sent in real time
func wsConnectionAPI(w http.ResponseWriter, req *http.Request) {
	// CORS handling
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade connection to websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// identifier for ws subscriber
	identifier, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// message channel for broker and subscribe
	msgCh := make(chan string)
	broker.AddSubscriber(identifier.String(), msgCh)

	// channel as signal to stop/close websocket
	stopCh := make(chan bool)

	// read from client
	go wsReadMsg(ws, stopCh)

	// write to client
	go wsWriteMsg(ws, identifier.String(), stopCh, msgCh)
}

func wsReadMsg(ws *websocket.Conn, stopCh chan<- bool) {
	_, _, errRead := ws.ReadMessage()

	// close connection when error occurred as closing closure will marked as error in GO
	if errRead != nil {
		stopCh <- true
		return
	}
}

func wsWriteMsg(ws *websocket.Conn, identifier string, stopCh <-chan bool, msgCh <-chan string) {
	// limit server to have 100msec delay before next processing to reduce heave load on server
	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		<-ticker.C

		// use non-blocking channel operation
		select {
		case <-stopCh:
			fmt.Printf("ID: %s. Stop signal. Closing websocket...\n", identifier)
			wsCloseConnection(identifier, ticker)
			return
		case msg := <-msgCh:
			fmt.Println("Get message from broker! Sending to client...")
			if errWrite := ws.WriteMessage(websocket.TextMessage, []byte(msg)); errWrite != nil {
				fmt.Printf("ID: %s. Error occurred. Closing websocket...\n", identifier)
				wsCloseConnection(identifier, ticker)
				return
			}
		default:
			fmt.Println("Waiting for message...")
		}
	}
}

func wsCloseConnection(identifier string, ticker *time.Ticker) {
	ticker.Stop()
	broker.RemoveSubscriber(identifier)
}