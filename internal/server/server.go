package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/magicznykacpur/httpfromtcp/internal/headers"
	"github.com/magicznykacpur/httpfromtcp/internal/request"
	"github.com/magicznykacpur/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("couldn't open listener on port %d: %v", port, err)
	}

	server := &Server{listener: listener, isClosed: atomic.Bool{}, handler: handler}
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

	resWriter := response.NewWriter()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		errorHeaders := headers.NewHeaders()
		errorHeaders.Set("Content-Type", "text/plain")

		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Body:       []byte(err.Error()),
			Headers:    errorHeaders,
		}
		hErr.WriteError(*resWriter)
		return
	}

	handlerErr := s.handler(*resWriter, req)
	if handlerErr != nil {
		err = handlerErr.WriteError(*resWriter)
		if err != nil {
			conn.Write([]byte(err.Error()))
			return
		}
		
		conn.Write(resWriter.Buffer.Bytes())
		return
	}

	conn.Write(resWriter.Buffer.Bytes())
}
