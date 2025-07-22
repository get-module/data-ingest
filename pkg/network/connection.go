package network

import (
	"bufio"
	"io"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	buffer := make([]byte, 4096) // adjustable buffer size per packet

	for {
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection closed by client")
			} else {
				log.Printf("read error: %v", err)
			}
			return
		}

		data := make([]byte, n)
		copy(data, buffer[:n])

		// TODO: send data to WAL
		log.Printf("received %d bytes: %x...", n, data[:min(16, n)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}