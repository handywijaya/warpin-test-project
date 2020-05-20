package pubsub

import (
	"fmt"
)

type Broker struct {
	Subscribers map[string]chan string
}

// add subscriber (websocket client)
func (br *Broker) AddSubscriber(identifier string, channel chan string) {
	_, found := br.Subscribers[identifier]
	if found {
		return
	}

	br.Subscribers[identifier] = channel
	fmt.Println("Subscriber ID: " + identifier + " registered!")
}

// remove subscriber (websocket client)
func (br *Broker) RemoveSubscriber(identifier string) {
	_, found := br.Subscribers[identifier]
	if !found {
		return
	}

	delete(br.Subscribers, identifier)

	fmt.Println("Subscriber ID: " + identifier + " unregistered!")
}

func (br *Broker) PublishMessage(msg string) {
	for _, ch := range br.Subscribers {
		ch <- msg
	}
}