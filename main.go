package main

import (
	"fmt"
	"os"
	"strings"
)

// *3 Indicates Array of size 3
// $5 Indicates String of size 5

func main() {
	config := buildConfig()
	conn, err := listen(config.Port)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	// Should I close the listener??

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if value.typ != "array" || len(value.array) == 0 {
			fmt.Println("Invalid request")
			continue
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		writer := NewWriter(conn)

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command", cmd)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		res := handler(args)
		writer.Write(res)
	}
}
