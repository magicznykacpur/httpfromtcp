package main

import (
	"fmt"
	"log"
	"net"

	"github.com/magicznykacpur/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Couldn't open tcp listener on port 42069: %v", err)
	}

	log.Println("listening for tcp request on :42069")

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Couldn't accept connection: %v", err)
		}

		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Printf("%v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
	}
}
