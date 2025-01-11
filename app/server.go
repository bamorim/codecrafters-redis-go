package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
)

const pong = "+PONG\r\n"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read data to
	b := make([]byte, 64)

	for {
		reader := bufio.NewReader(conn)

		value, err := Parse(reader)

		fmt.Printf("Received Value: %#v\n", value)

		if err != nil {
			fmt.Println("Error reading bytes from connection")
			// Assume connection was closed and just continue
			break
		}

		// Right now we are ignoring everything and just waiting until the end of line
		if bytes.Contains(b, []byte("\n")) {
			conn.Write([]byte(pong))
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
