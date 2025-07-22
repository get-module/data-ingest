package parser

import (
	"fmt"
	"sync"
	"error"
	"pkg/config/config"
)

type Message struct{
	source    string
	timestamp int
	message   map[string]interface{}
}

func parse(data []bytes){
	var errOnce sync.Once
	var firstErr error
	var wg sync.WaitGroup
	

	/* verification checks on the data */
	go verify_length(&data, &wg)
	go verify_json(&data, &wg)
	

	wg.Wait() // ends function once all routines are done

	if firstErr != nil {
		return firstErr.Error()
	}

	// parsing logic

	return Message{}
}

func verify_length(data *[]byte, wg *sync.WaitGroup){
	defer wg.Done()

	if (len(data) != int(config.env("PACKET_SIZE"))){
		errOnce.Do(func(){ firstErr = "Invalid Packet Size - " + string(len(data)) })
	}
}

func verify_json(data *[]byte){
	defer wg.Done()

	var js json.RawMessage
	if (!json.Unmarshal(data, &js)){
		errOnce.Do(func(){ firstErr = "Invalid JSON format" })
		return
	}

	keys = ["source", "timestamp", "message"]

	for _, key := range keys{
		if _, ok := obj[key]; !ok{
			errOnce.Do(func(){ firstErr = "Missing Key - " + string(key) })
		}
	}
}

// maybe have logic where if the len() != the preset definition, the proper len is sent to the client to update (if different for whatever reason) to minimize breaking?
