package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/a98c14/hyperion/router"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleGameConnection() {

}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4000"
	CONN_TYPE = "tcp"
)

func handleRequest(conn net.Conn) {
	fmt.Println("Received message...")
	buffer := make([]byte, 1024)
	byteCount, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Could not read message", err.Error())
		return
	}

	fmt.Println("Read " + strconv.Itoa(byteCount) + " bytes")
	fmt.Println(string(buffer))
	conn.Write([]byte("Message received"))
	conn.Close()
}

func listenGameSocket() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on: " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func startWebServer() {
	fmt.Println("Starting web server...")
	http.ListenAndServe("127.0.0.1:8000", router.New())
}

func main() {
	// TODO(selim):
	// - Listen for websocket connections and http requests (For custom editor)
	// - Create a structure for currently active games
	// - Some sort of communication between server and active games (be able to run commands etc.)
	// - Database connection (local for now, maybe sqlite?)
	// - Endpoint for loading latest balance values.
	// - Endpoint for loading all balance value list.
	// - Endpoint for saving current balance values.
	fmt.Println("App started.")
	go listenGameSocket()
	startWebServer()
}
