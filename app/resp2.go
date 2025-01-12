package main

type ValueType int

//go:generate stringer -type=ValueType
const (
	// NullBulkString should be the zero value to make it easier to deal with nulls
	NullBulkStringType ValueType = iota
	SimpleStringType
	SimpleErrorType
	IntegerType
	BulkStringType
	NullArrayType
	ArrayType
)

type Value struct {
	Type    ValueType
	Integer int64
	String  string
	Values  []Value
}

func NewInteger(i int64) Value {
	return Value{
		Type:    IntegerType,
		Integer: i,
	}
}

func NewSimpleString(s string) Value {
	return Value{
		Type:   SimpleStringType,
		String: s,
	}
}

func NewSimpleError(s string) Value {
	return Value{
		Type:   SimpleErrorType,
		String: s,
	}
}

func NewBulkString(s string) Value {
	return Value{
		Type:   BulkStringType,
		String: s,
	}
}

func NewNullBulkString() Value {
	return Value{
		Type: NullBulkStringType,
	}
}

func NewArray(values []Value) Value {
	return Value{
		Type:   ArrayType,
		Values: values,
	}
}

func NewNullArray() Value {
	return Value{
		Type: NullArrayType,
	}
}
