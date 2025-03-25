package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	line := ""

	go func() {
		for {
			bytes := make([]byte, 8)

			_, err := f.Read(bytes)
			if err == io.EOF {
				if line != "" {
					lines <- line
				}

				f.Close()
				close(lines)

				return
			}

			parts := strings.Split(string(bytes), "\n")

			if len(parts) == 1 {
				line += parts[0]
			} else {
				line += parts[0]
				lines <- line

				line = ""
				line += parts[1]
			}
		}
	}()

	return lines
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Couldn't open tcp listener on port 42069: %v", err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Couldn't accept connection: %v", err)
		}

		lines := getLinesChannel(connection)

		for line := range lines {
			fmt.Println(line)
		}
	}
}
