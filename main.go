package main

import (
	"fmt"
	"log"

	tcpServer "bin.go.mod/tcp"
)

func main() {
	server := tcpServer.NewServer(":8000")

	go func() {

		if err := server.Start(); err != nil {
			log.Fatal(err)
		}

	}()

	for msg := range server.Msgs {
		fmt.Printf("message from (%s):%s", msg.From, string(msg.Payload))
	}
}
