package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type Memory struct {
	l  sync.Mutex
	kv map[string]Value
}

func (memory *Memory) Init() {
	memory.kv = make(map[string]Value)
}

func (memory *Memory) Set(key string, value Value) {
	memory.l.Lock()
	defer memory.l.Unlock()

	memory.kv[key] = value
}

func (memory *Memory) Get(key string) Value {
	return memory.kv[key]
}

var memory Memory

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

	response := responseFor(command, args)

	if error := Write(w, response); error != nil {
		return error
	}

	return nil
}

func responseFor(command string, args []string) Value {
	switch command {
	case "PING":
		return NewSimpleString("PONG")
	case "ECHO":
		if len(args) != 1 {
			return NewSimpleError("ERR ECHO expects exactly one argument")
		}
		return NewBulkString(args[0])
	case "SET":
		if len(args) < 2 {
			return NewSimpleError("ERR SET requires key and value")
		}

		memory.Set(args[0], NewBulkString(args[1]))

		return NewSimpleString("OK")
	case "GET":
		if len(args) < 1 {
			return NewSimpleError("ERR GET requires key")
		}

		return memory.Get(args[0])
	}
	return NewSimpleError("ERR undefined command")
}

func normalizeCommand(message Value) (string, []string, error) {
	if message.Type != ArrayType {
		return "", []string{}, fmt.Errorf("command is not an array")
	}

	if len(message.Values) < 1 {
		return "", []string{}, fmt.Errorf("command array is empty")
	}

	for _, value := range message.Values {
		if value.Type != BulkStringType {
			return "", []string{}, fmt.Errorf("command argument is not a BulkString")
		}
	}

	arguments := make([]string, len(message.Values)-1)

	for i, value := range message.Values[1:] {
		arguments[i] = value.String
	}

	return message.Values[0].String, arguments, nil
}

func main() {
	// Initialize the shared memory
	memory.Init()

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
