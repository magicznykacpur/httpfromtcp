package response

import (
	"io"
	"strconv"

	"github.com/magicznykacpur/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok                  = 200
	BadRequest          = 400
	InternalServerError = 500
)

const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case Ok:
		_, err := w.Write([]byte("HTTP/1.1 200 OK" + crlf))
		if err != nil {
			return err
		}
	case BadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request" + crlf))
		if err != nil {
			return err
		}
	case InternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error" + crlf))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()

	defaultHeaders.Set("Content-Length", strconv.Itoa(contentLen))
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")

	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + crlf))
		if err != nil {
			return err
		}
	}
	return nil
}
