package main

import (
	"fmt"
	"os"
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

		fmt.Printf("Received value: %v\n", value)
		conn.Write([]byte("+OK\r\n"))
	}
}
