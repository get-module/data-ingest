package main

import (
	"log"
	"github.com/get-module/data-ingest/pkg/network"
)

func main(){
	err := network.Start(":8080", network.HandleConnection)
	if err != nil{
		log.Fatalf("Server error: %v", err)
	}
}
