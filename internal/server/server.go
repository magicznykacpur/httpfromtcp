package server

import (
	"bytes"
	"fmt"
	"net"
	"sync/atomic"

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		hErr.WriteError(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		handlerErr.WriteError(conn)
		return
	}

	defaultHeaders := response.GetDefaultHeaders(len(buf.Bytes()))

	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.InternalServerError,
			Message:    err.Error(),
		}
		hErr.WriteError(conn)
		return
	}

	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.InternalServerError,
			Message:    err.Error(),
		}
		hErr.WriteError(conn)
		return
	}

	conn.Write(buf.Bytes())
}
