package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	ExpiresAt time.Time
	// Use negated value to exploit zero-value so a zero value entry is expired
	Infinite bool
	Value    Value
}

type Memory struct {
	Lock sync.Mutex
	KV   map[string]Entry
}

func (memory *Memory) Init() {
	memory.KV = make(map[string]Entry)
}

func (memory *Memory) Set(key string, entry Entry) {
	memory.Lock.Lock()
	defer memory.Lock.Unlock()

	memory.KV[key] = entry
}

func (memory *Memory) Get(key string) Value {
	entry := memory.KV[key]

	if entry.Infinite || entry.ExpiresAt.After(time.Now()) {
		return memory.KV[key].Value
	} else {
		return NewNullBulkString()
	}
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
	switch strings.ToUpper(command) {
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

		// Required args
		key := args[0]
		entry := Entry{Value: NewBulkString(args[1]), Infinite: true}

		// Parse optional args
		for i := 2; i < len(args); i++ {
			switch strings.ToUpper(args[i]) {
			case "PX":
				if len(args) < i+2 {
					return NewSimpleError("ERR SET PX option requires an argument")
				}

				px, error := strconv.ParseInt(args[i+1], 10, 0)

				if error != nil {
					return NewSimpleError("ERR SET PX option argument must be a positive integer")
				}

				entry.Infinite = false
				entry.ExpiresAt = time.Now().Add(time.Millisecond * time.Duration(px))
			}
		}

		memory.Set(key, entry)

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
