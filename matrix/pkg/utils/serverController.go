/*
Copyright © [2022] [Lakshy Sharma] <lakshy.sharma@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

/*
TCP server functions.
*/
func processTCPClient(clientConnection net.Conn, replyMessage string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Received a new client connection.")
	for {
		message, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		// Exit when client sends EOF
		if message == "EOF" {
			return
		}

		// Display the received message.
		fmt.Println("<- ", string(message))

		// Client asked for a echo server then reformat the message and send back.
		// Else simply send back the required response.
		if replyMessage == "ECHO" {
			clientConnection.Write([]byte(fmt.Sprintf("Echo: %s", string(message))))
		} else {
			clientConnection.Write([]byte(string(replyMessage + "\n")))
		}
	}
}

// This function starts a TCP server with provided port and reply mechanism.
func ServeTCP(portNumber int, replyMessage string) {
	wg := new(sync.WaitGroup)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(portNumber))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("TCP server started.\nPort: %d\nReply: %s\n", portNumber, replyMessage)

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go processTCPClient(clientConnection, replyMessage, wg)
	}
}

/*
Websocket server functions.
*/

var upgrader = websocket.Upgrader{}
var replyMessage string

func socketHandler(w http.ResponseWriter, r *http.Request) {
	websocketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during upgrading connection: ", err)
		return
	}
	defer websocketConnection.Close()

	for {
		messageType, message, err := websocketConnection.ReadMessage()
		if err != nil {
			if fmt.Sprintf("%s", err) == "websocket: close 1006 (abnormal closure): unexpected EOF" {
				log.Println("Client Exited.")
				break
			}
			log.Println("Error during reading client message: ", err)
			break
		}
		fmt.Printf("<- %s", message)

		// Writing message back to the client.
		if replyMessage == "ECHO" {
			reply := "Echo: " + string(message)
			err = websocketConnection.WriteMessage(messageType, []byte(reply))
			if err != nil {
				log.Println("Error during writing message to socket: ", err)
				break
			}
		} else {
			reply := replyMessage
			err = websocketConnection.WriteMessage(messageType, []byte(reply))
			if err != nil {
				log.Println("Error during writing message to socket: ", err)
				break
			}
		}

	}
}

// This function starts a Websocket server.
func ServeWebsocket(portNumber int, replyString string) {
	replyMessage = replyString
	http.HandleFunc("/", socketHandler)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(portNumber), nil))
}
