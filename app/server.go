package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const pong = "+PONG\r\n"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)

		message, err := Parse(reader)
		if err != nil {
			fmt.Println("Error reading message from connection")
			// Assume connection was closed and just continue
			// TODO: Probably we should fallback to the telnet protocol version (inline commands)
			break
		}

		if error := processCommand(conn, message); error != nil {
			fmt.Println("Could not process command")
			break
		}
	}
}

func processCommand(w io.Writer, message Value) error {
	command, args, error := normalizeCommand(message)

	if error != nil {
		return error
	}

	switch command {
	case "PING":
		Write(w, NewSimpleString("PONG"))
	case "ECHO":
		Write(w, NewArray(args))
	default:
		Write(w, NewSimpleError("ERR undefined command"))
	}

	return nil
}

func normalizeCommand(message Value) (string, []Value, error) {
	if message.Type != ArrayType {
		return "", []Value{}, fmt.Errorf("command is not an array")
	}

	if len(message.Values) < 1 {
		return "", []Value{}, fmt.Errorf("command array is empty")
	}

	for _, value := range message.Values {
		if value.Type != BulkStringType {
			return "", []Value{}, fmt.Errorf("command argument is not a BulkString")
		}
	}

	return message.Values[0].String, message.Values[1:], nil
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
