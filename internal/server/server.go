package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/magicznykacpur/httpfromtcp/internal/response"
)

type Server struct {
	port     int
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("couldn't open listener on port %d: %v", port, err)
	}

	server := &Server{port: port, listener: listener, isClosed: atomic.Bool{}}
	server.isClosed.Store(false)

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("couldn't close listener: %v", err)
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			fmt.Println("couldn't accept connection")
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		fmt.Println(err)
	}
	
	defaultHeaders := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		fmt.Println(err)
	}
}
