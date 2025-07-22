package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// Packet represents the JSON packet structure
type Packet struct {
	Checksum  string `json:"checksum"`
	Source    string `json:"source"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	const serverAddr = "127.0.0.1:8080"
	const packetsPerSecond = 10
	interval := time.Second / packetsPerSecond

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; ; i++ {
		<-ticker.C

		dataStr := randomString(100) // random 100 char string as data

		// Calculate checksum (SHA256 of data)
		hash := sha256.Sum256([]byte(dataStr))
		checksum := hex.EncodeToString(hash[:])

		pkt := Packet{
			Checksum:  checksum,
			Source:    fmt.Sprintf("client-%d", rand.Intn(1000)),
			Data:      dataStr,
			Timestamp: time.Now().Format(time.RFC3339Nano),
		}

		jsonBytes, err := json.Marshal(pkt)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			continue
		}

		// Send JSON + newline (optional, if your server expects delimiter)
		_, err = conn.Write(append(jsonBytes, '\n'))
		if err != nil {
			log.Printf("Failed to send packet: %v", err)
			break
		}

		fmt.Printf("Sent packet #%d: %s\n", i+1, string(jsonBytes))
	}
}
