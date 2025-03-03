package main

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
