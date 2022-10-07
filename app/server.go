package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(2)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		if _, err := conn.Read([]byte{}); err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}
		_, err := conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(3)
		}
	}
}
