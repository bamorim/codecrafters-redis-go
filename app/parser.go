package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func expect(reader *bufio.Reader, byte byte) error {
	next, error := reader.ReadByte()

	if error != nil {
		// TODO: Treat error better
		return error
	}

	if next != byte {
		return fmt.Errorf("expected %c found %c", byte, next)
	}

	return nil
}

func parseSimple(reader *bufio.Reader, byte byte) (string, error) {
	if error := expect(reader, byte); error != nil {
		return "", error
	}

	string, error := reader.ReadString('\r')

	if error != nil {
		// TODO: Treat error better
		return "", error
	}

	if error := expect(reader, '\n'); error != nil {
		// TODO: Treat error better
		return "", error
	}

	return strings.TrimRight(string, "\r"), nil
}

func parseSimpleString(reader *bufio.Reader) (Value, error) {
	string, error := parseSimple(reader, '+')
	if error != nil {
		return Value{}, error
	}

	return NewSimpleString(string), nil
}

func parseErrorString(reader *bufio.Reader) (Value, error) {
	string, error := parseSimple(reader, '-')
	if error != nil {
		return Value{}, error
	}

	return NewSimpleError(string), nil
}

func parseInteger(reader *bufio.Reader) (Value, error) {
	string, error := parseSimple(reader, ':')
	if error != nil {
		return Value{}, error
	}

	integer, error := strconv.ParseInt(string, 10, 0)
	if error != nil {
		return Value{}, error
	}

	return NewInteger(integer), nil
}

func parseBulkString(reader *bufio.Reader) (Value, error) {
	lengthString, error := parseSimple(reader, '$')
	if error != nil {
		return Value{}, error
	}

	length, error := strconv.Atoi(lengthString)
	if error != nil {
		return Value{}, error
	}

	if length < 0 {
		return NewNullBulkString(), nil
	}

	buff := make([]byte, length)

	read, error := reader.Read(buff)

	if read != length {
		return Value{}, fmt.Errorf("Expected a string with %d bytes but only read %d", length, read)
	}

	expect(reader, '\r')
	expect(reader, '\n')

	return NewBulkString(string(buff)), nil
}

func parseArray(reader *bufio.Reader) (Value, error) {
	lengthString, error := parseSimple(reader, '*')
	if error != nil {
		return Value{}, error
	}

	length, error := strconv.Atoi(lengthString)
	if error != nil {
		return Value{}, error
	}

	if length < 0 {
		return NewNullArray(), nil
	}

	values := make([]Value, length)

	for i := 0; i < length; i++ {
		values[i], error = parseValue(reader)

		if error != nil {
			return Value{}, error
		}
	}

	return NewArray(values), nil
}

// TODO: Avoid recursion stack overflow somehow
func parseValue(reader *bufio.Reader) (Value, error) {
	next, error := reader.Peek(1)
	if error != nil {
		return Value{}, error
	}

	switch next[0] {
	case '+':
		return parseSimpleString(reader)
	case '-':
		return parseErrorString(reader)
	case ':':
		return parseInteger(reader)
	case '$':
		return parseBulkString(reader)
	case '*':
		return parseArray(reader)
	}

	return Value{}, fmt.Errorf("invalid value prefix: %c", next[0])
}

func Parse(reader *bufio.Reader) (Value, error) {
	return parseValue(reader)
}
