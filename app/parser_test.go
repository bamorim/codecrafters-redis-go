package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseSimpleString(t *testing.T) {
	input := strings.NewReader("+PING\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, SimpleStringType)
	assertString(t, result, "PING")
}

func TestParseSimpleError(t *testing.T) {
	input := strings.NewReader("-ERR something wrong\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, SimpleErrorType)
	assertString(t, result, "ERR something wrong")
}

func TestParseInteger(t *testing.T) {
	input := strings.NewReader(":1234\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, IntegerType)
	assertInteger(t, result, 1234)
}

func TestParseBulkString(t *testing.T) {
	input := strings.NewReader("$10\r\nhello\r\n123\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, BulkStringType)
	assertString(t, result, "hello\r\n123")
}

func TestParseNullBulkString(t *testing.T) {
	input := strings.NewReader("$-1\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, NullBulkStringType)
}

func TestParseEmptyBulkString(t *testing.T) {
	input := strings.NewReader("$0\r\n\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, BulkStringType)
	assertString(t, result, "")
}

func TestParseArray(t *testing.T) {
	input := strings.NewReader("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, ArrayType)

	if len(result.Values) != 2 {
		t.Fatalf("expected array to have 2 elements but instead got %d elements", len(result.Values))
	}

	assertType(t, result.Values[0], BulkStringType)
	assertString(t, result.Values[0], "ECHO")
	assertType(t, result.Values[1], BulkStringType)
	assertString(t, result.Values[1], "hey")
}

func TestParseNullArray(t *testing.T) {
	input := strings.NewReader("*-1\r\n")
	reader := bufio.NewReader(input)
	result, error := Parse(reader)

	refuteError(t, error)
	assertType(t, result, NullArrayType)
}

func TestParseEmptyInputFail(t *testing.T) {
	input := strings.NewReader("")
	reader := bufio.NewReader(input)
	_, error := Parse(reader)

	assertError(t, error)
}

func refuteError(t *testing.T, error error) {
	if error != nil {
		t.Fatalf("Expected parse to not return error, got %#v", error)
	}
}

func assertError(t *testing.T, error error) {
	if error == nil {
		t.Fatal("Expected parse to return error")
	}
}

func assertType(t *testing.T, result Value, vt ValueType) {
	if result.Type != vt {
		t.Fatalf("Expected result to be of type %s, got %s", vt, result.Type)
	}
}

func assertString(t *testing.T, result Value, string string) {
	if result.String != string {
		t.Fatalf("Expected result value to be the string %#v, got %#v", string, result.String)
	}
}

func assertInteger(t *testing.T, result Value, integer int64) {
	if result.Integer != integer {
		t.Fatalf("Expected result value to be the integer %#v, got %#v", integer, result.Integer)
	}
}
