package main

import (
	"strconv"
)

type Value struct {
	// datatype
	typ string // cannot use `type` as it's a reserved keyword
	// simple string value
	str string
	// numerical value
	num int
	// bulk string value
	bulk string
	// array value
	array []Value
}

const (
	CLRF = "\r\n"
)

func (v Value) marshalArray() []byte {
	len := len(v.array)
	var b []byte
	b = append(b, ARRAY)
	b = append(b, strconv.Itoa(len)...)
	b = append(b, CLRF...)

	for i := range len {
		b = append(b, v.array[i].Marshal()...)
	}

	return b
}

func (v Value) marshalString() []byte {
	var b []byte
	b = append(b, STRING)
	b = append(b, v.str...)
	b = append(b, CLRF...)

	return b
}

func (v Value) marshalBulk() []byte {
	var b []byte
	b = append(b, BULK_STRING)
	b = append(b, strconv.Itoa(len(v.bulk))...)
	b = append(b, CLRF...)

	b = append(b, v.bulk...)
	b = append(b, CLRF...)

	return b
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalError() []byte {
	var b []byte
	b = append(b, ERROR)
	b = append(b, v.str...)
	b = append(b, CLRF...)

	return b
}

// Transforms message body into a binary / textual format
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}
