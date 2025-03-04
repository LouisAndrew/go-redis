package main

import (
	"fmt"
	"net"
)

func listen(port int) (net.Listener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening on port:", port)

	return l, nil
}
