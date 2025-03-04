package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	config := buildConfig()

	db, err := NewAof(config.AofPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db.Read(func(v Value) {
		handleValue(v)
	})

	l, err := listen(config.Port)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go handleConnection(conn, db)
	}
}

func handleValue(v Value) (Value, string) {
	cmd := strings.ToUpper(v.array[0].bulk)
	args := v.array[1:]
	if v.typ != "array" || len(v.array) == 0 {
		fmt.Println("Invalid request")
		return Value{typ: "error", str: "Invalid request"}, cmd
	}

	handler, ok := Handlers[cmd]
	if !ok {
		return Value{typ: "error", str: "Invalid command"}, cmd
	}

	return handler(args), cmd
}

func handleConnection(conn net.Conn, db *Aof) {
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		writer := NewWriter(conn)
		res, cmd := handleValue(value)

		if cmd == "SET" || cmd == "HSET" {
			db.Write(value)
		}

		writer.Write(res)
	}
}
