package server

import (
	"fmt"
	"io"

	"github.com/magicznykacpur/httpfromtcp/internal/request"
	"github.com/magicznykacpur/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(io.Writer, *request.Request) *HandlerError

func (hr *HandlerError) WriteError(w io.Writer) error {
	err := response.WriteStatusLine(w, hr.StatusCode)
	if err != nil {
		fmt.Println("could write status line for handler error")
	}
	w.Write([]byte("\r\n" + hr.Message))

	return nil
}
