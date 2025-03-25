package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Couldn't resolve UDP address: %v", err)
	}

	connection, err := net.DialUDP("udp", nil, listener)
	if err != nil {
		log.Fatalf("Couldnt dial a UDP connection: %v", err)
	}

	defer connection.Close()

	for {
		fmt.Printf(">")
		reader := bufio.NewReader(os.Stdin)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Couldn't read line: %v", err)
		}

		_, err = connection.Write([]byte(line))
		if err != nil {
			fmt.Printf("Couldn't write line: %v", err)
		}
	}
}