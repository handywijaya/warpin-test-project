# How to use
Run this code by using
```
go run src/main/main.go
```
and visit `localhost:8090`

# API List
There are several endpoints that can be used
1. `/` is for home endpoint
1. `/send` to send messages. The message will be saved in a file named `data.txt`
1. `/get` to get all sent messages.
1. `/flush` to clean all sent messages.
1. `/ws` to open an HTML page containing button to connect to server using websocket. Try to send a message using `/send` endpoint in a new tab and witness the magic.

# Credits
This project is using
1. [`github.com/gorilla/websocket`](https://github.com/gorilla/websocket) to implement websocket.
1. [`github.com/nu7hatch/gouuid`](https://github.com/nu7hatch/gouuid) to get generate a unique identifier for subscribers.