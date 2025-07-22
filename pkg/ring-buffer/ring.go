package ringbuffer

import (
	"fmt"
	"github.com/smallnest/ringbuffer"
	"pkg/config/config"
)

// TODO: replace all instances of 4096 with the CONFIG file read
// TODO: implement error handling

type RingBuffer struct {
	buffer [][]byte
	size   int
	start  int
	end    int
	full   bool
}

type RingBufferReturn struct {
	length int
	free   int
}

func NewRingBuffer(capacity int) *RingBuffer{
	// capacity is the number of pieces of data it can contain (RB_MAX_SUBSCRIBERS)
	capacity = capacity * int(config.env("PACKET_SIZE"))

	return ringbuffer.New(capacity).setBlocking(false) // blocking allows for overwrites
}

func (rb *RingBuffer) Push(data []byte) RingBufferReturn{
	// TODO: sizing logic (4096 bytes etc etc)
	
	rb.Write([]byte(data))

	return RingBufferReturn{length: rb.Length(), free: rb.Free()}
}

func (rb *RingBuffer) Read(indexes int){
	buffer := make([]byte, int(config.env("PACKET_SIZE")) * indexes) // 4096 bytes is each piece of data here	
	rb.Read(buffer)
	return string(buffer)
}
