package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const pong = "+PONG\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	// Create a buffer to read data to
	b := make([]byte, 64)

	for {
		bc, err := conn.Read(b)
		fmt.Println("Read some bytes", bc, string(b))
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
