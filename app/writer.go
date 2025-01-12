package main

import (
	"fmt"
	"io"
	"strconv"
)

func Write(w io.Writer, message Value) error {
	suffix, error := suffixFor(message)

	if error != nil {
		return error
	}

	if _, error := w.Write([]byte{suffix}); error != nil {
		return error
	}

	switch message.Type {
	case SimpleStringType, SimpleErrorType:
		if error := writeLine(w, message.String); error != nil {
			return error
		}
	case IntegerType:
		if error := writeLine(w, strconv.FormatInt(message.Integer, 10)); error != nil {
			return error
		}
	case BulkStringType:
		if error := writeLine(w, strconv.Itoa(len(message.String))); error != nil {
			return error
		}

		if error := writeLine(w, message.String); error != nil {
			return error
		}
	case ArrayType:
		if error := writeLine(w, strconv.Itoa(len(message.Values))); error != nil {
			return error
		}

		for _, value := range message.Values {
			if error := Write(w, value); error != nil {
				return error
			}
		}

	case NullBulkStringType, NullArrayType:
		if error := writeLine(w, "-1"); error != nil {
			return error
		}
	}

	return nil
}

func suffixFor(message Value) (byte, error) {
	switch message.Type {
	case SimpleStringType:
		return '+', nil

	case SimpleErrorType:
		return '-', nil

	case IntegerType:
		return ':', nil

	case BulkStringType, NullBulkStringType:
		return '$', nil
	case ArrayType, NullArrayType:
		return '*', nil
	}

	return 0, fmt.Errorf("invalid suffix for message value type %v", message.Type)
}

func writeLine(w io.Writer, s string) error {
	if _, error := w.Write([]byte(s)); error != nil {
		return error
	}
	if _, error := w.Write([]byte{'\r', '\n'}); error != nil {
		return error
	}

	return nil
}
