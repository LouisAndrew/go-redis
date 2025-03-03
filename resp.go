package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING      = '+' // Simple string, like `+OK`. Never contains `\r` and `\n`
	ERROR       = '-' // Simple error message
	INTEGER     = ':'
	BULK_STRING = '$' // Binary string. Can be encoded as a blob. Use this for usual strings
	ARRAY       = '*'
)

type Resp struct {
	reader *bufio.Reader // Using `bufio` so that we can use the buffering functionalities
}

// Returns the line, number of bytes and a potential error
// the `reader` has an internal state (e.g. whenever you do `readByte`, it will)
// store the last position.
// Returning a copy of Resp -> `r` will be copied over and over.
// Means, the pointer will also be copied. This might not be a good idea.
// Returning pointer means you're using the same `r` instance over all operations
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		// Below code you don't want to do!
		// Why: You want to make sure that the line is fully consumed by the reader,
		// until the end (\r\n), since the `reader` stores internal position of where
		// Exactly the buffer was last read
		// if b == '\r' {
		// 	break
		// }

		n += 1
		line = append(line, b)

		// Line is fully consumed. Check for linebreaks here
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	// Will this still be safe if length of `line` is less than 2?
	return line[:len(line)-2], n, nil
}

// Returns the integer of a line, number of bytes and potential error
// See `readLine` docs for more info
func (r *Resp) readInt() (val int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(i64), n, nil
}

func (r *Resp) readArray() (val Value, err error) {
	val.typ = "array"
	length, _, err := r.readInt()
	if err != nil {
		return val, err
	}

	val.array = make([]Value, length)
	for i := 0; i < length; i++ {
		el, err := r.Read()
		if err != nil {
			return val, err
		}

		val.array[i] = el
	}

	return val, nil
}

func (r *Resp) readBulkString() (val Value, err error) {
	val.typ = "bulk"
	length, _, err := r.readInt()
	if err != nil {
		return val, err
	}

	b := make([]byte, length)
	r.reader.Read(b)
	val.bulk = string(b)
	r.readLine() // Read trailing CRLF (`\r\n`)
	return val, nil
}

// Public method to read from connection input
func (r *Resp) Read() (val Value, err error) {
	dataType, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case ARRAY:
		return r.readArray()
	case BULK_STRING:
		return r.readBulkString()
	default:
		return Value{}, fmt.Errorf("Unknown type %s", string(dataType))
	}
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{bufio.NewReader(reader)}
}
