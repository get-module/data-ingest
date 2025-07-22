package network

import (
	"fmt"
	"net"
	"log"
)

func Start(addr string, handler func(conn net.Conn)) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil{
		return fmt.Errorf("Failed to start server: %w", err)
	}

	log.Printf("TCP server listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		SetSocketOptions(conn) 
		go handler(conn)
	}
}