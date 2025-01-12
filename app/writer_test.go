package main

import (
	"strings"
	"testing"
)

func TestWriteSimpleString(t *testing.T) {
	message := NewSimpleString("PING")
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "+PING\r\n")
}

func TestWriteSimpleError(t *testing.T) {
	message := NewSimpleError("ERR something")
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "-ERR something\r\n")
}

func TestWriteInteger(t *testing.T) {
	message := NewInteger(12345)
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), ":12345\r\n")
}

func TestWriteBulkString(t *testing.T) {
	message := NewBulkString("Hello\r\nWorld")
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "$12\r\nHello\r\nWorld\r\n")
}

func TestWriteNullBulkString(t *testing.T) {
	message := NewNullBulkString()
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "$-1\r\n")
}

func TestWriteEmptyBulkString(t *testing.T) {
	message := NewBulkString("")
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "$0\r\n\r\n")
}

func TestWriteArray(t *testing.T) {
	message := NewArray([]Value{
		NewBulkString("ECHO"),
		NewBulkString("hey"),
	})
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")
}

func TestWriteNullArray(t *testing.T) {
	message := NewNullArray()
	var builder strings.Builder
	error := Write(&builder, message)
	refuteError(t, error)
	assertWriteResult(t, builder.String(), "*-1\r\n")
}

func TestFailWithMalformedValue(t *testing.T) {
	message := Value{Type: 9999}
	var builder strings.Builder
	error := Write(&builder, message)
	assertError(t, error)
}

func assertWriteResult(t *testing.T, result string, expected string) {
	if result != expected {
		t.Fatalf("Expected writer to write the string %#v, but got instead %#v", expected, result)
	}
}
