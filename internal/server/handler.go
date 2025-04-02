package server

import (
	"fmt"

	"github.com/magicznykacpur/httpfromtcp/internal/headers"
	"github.com/magicznykacpur/httpfromtcp/internal/request"
	"github.com/magicznykacpur/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Headers    headers.Headers
	Body       []byte
}

type Handler func(response.Writer, *request.Request) *HandlerError

func (hr *HandlerError) WriteError(w response.Writer) error {
	err := w.WriteStatusLine(hr.StatusCode)
	if err != nil {
		return fmt.Errorf("couldn't write status line for handler error")
	}

	err = w.WriteHeaders(hr.Headers)
	if err != nil {
		return fmt.Errorf("couldn't write headers for handler error")
	}

	_, err = w.WriteBody(hr.Body)
	if err != nil {
		return fmt.Errorf("couldn't write body for handler error")
	}

	return nil
}
