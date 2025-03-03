package main

import (
	"fmt"
	"net"
)

// @TODO: Concurrent connections -> https://opensource.com/article/18/5/building-concurrent-tcp-server-go
func listen(port int) (conn net.Conn, err error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening on port:", port)

	// accept next connection to port
	_conn, err := l.Accept()
	if err != nil {
		return nil, err
	}

	return _conn, nil
}
