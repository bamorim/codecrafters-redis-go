// Code generated by "stringer -type=ValueType"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SimpleStringType-0]
	_ = x[SimpleErrorType-1]
	_ = x[IntegerType-2]
	_ = x[NullBulkStringType-3]
	_ = x[BulkStringType-4]
	_ = x[NullArrayType-5]
	_ = x[ArrayType-6]
}

const _ValueType_name = "SimpleStringTypeSimpleErrorTypeIntegerTypeNullBulkStringTypeBulkStringTypeNullArrayTypeArrayType"

var _ValueType_index = [...]uint8{0, 16, 31, 42, 60, 74, 87, 96}

func (i ValueType) String() string {
	if i < 0 || i >= ValueType(len(_ValueType_index)-1) {
		return "ValueType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ValueType_name[_ValueType_index[i]:_ValueType_index[i+1]]
}
