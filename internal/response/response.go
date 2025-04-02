package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/magicznykacpur/httpfromtcp/internal/headers"
)

type Writer struct {
	Buffer *bytes.Buffer
	state  writerState
}

type writerState int

const (
	writerStateInitialized writerState = iota
	writerStateWritingHeaders
	writerStateWritingBody
)

type StatusCode int

const (
	StatusOk                  = 200
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

const crlf = "\r\n"
const protocol = "HTTP/1.1"

func NewWriter() *Writer {
	return &Writer{Buffer: bytes.NewBuffer([]byte{}), state: writerStateInitialized}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var statusLine string

	switch statusCode {
	case StatusOk:
		statusLine = fmt.Sprintf("%s %d %s %s", protocol, StatusOk, "OK", crlf)
	case StatusBadRequest:
		statusLine = fmt.Sprintf("%s %d %s %s", protocol, StatusBadRequest, "Bad Request", crlf)
	case StatusInternalServerError:
		statusLine = fmt.Sprintf("%s %d %s %s", protocol, StatusInternalServerError, "Internal Server Error", crlf)
	default:
		return fmt.Errorf("unknown status code")
	}

	_, err := w.Buffer.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.state = writerStateWritingHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateWritingHeaders {
		return fmt.Errorf("invalid writer status, write status line first")
	}

	for key, value := range headers {
		_, err := w.Buffer.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
		if err != nil {
			return err
		}
	}

	_, err := w.Buffer.Write([]byte(crlf))
	if err != nil {
		return err
	}

	w.state = writerStateWritingBody

	return nil
}

func (w *Writer) WriteBody(bytes []byte) (int, error) {
	if w.state != writerStateWritingBody {
		return 0, fmt.Errorf("invalid writer status, write headers first")
	}

	return w.Buffer.Write(bytes)
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
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
