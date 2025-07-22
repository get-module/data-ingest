package ringbuffer

import (
	"fmt"
	"github.com/smallnest/ringbuffer"
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
	return ringbuffer.New(capacity).setBlocking(false) // blocking allows for overwrites
}

func (rb *RingBuffer) Push(data []byte) RingBufferReturn{
	// TODO: sizing logic (4096 bytes etc etc)
	rb.Write([]byte(data))

	return RingBufferReturn{length: rb.Length(), free: rb.Free()}
}

func (rb *RingBuffer) Read(indexes int){
	buffer := make([]byte, 4096 * indexes) // 4096 bytes is each piece of data here	
	rb.Read(buffer)
	return string(buffer)
}
