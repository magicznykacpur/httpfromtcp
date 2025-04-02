package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/magicznykacpur/httpfromtcp/internal/headers"
	"github.com/magicznykacpur/httpfromtcp/internal/request"
	"github.com/magicznykacpur/httpfromtcp/internal/response"
	"github.com/magicznykacpur/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w response.Writer, r *request.Request) *server.HandlerError {
	errorHeaders := headers.NewHeaders()
	errorHeaders.Set("Content-Type", "text/html")

	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Headers:    errorHeaders,
			Body:       badRequest,
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Headers:    errorHeaders,
			Body:       internalServerError,
		}
	default:
		err := w.WriteStatusLine(response.StatusOk)
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Headers:    errorHeaders,
				Body:       internalServerError,
			}
		}

		headers := response.GetDefaultHeaders(len(okRequest))
		headers.OverrideHeader("Content-Type", "text/html")

		err = w.WriteHeaders(headers)
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Headers:    errorHeaders,
				Body:       internalServerError,
			}
		}

		_, err = w.WriteBody(okRequest)
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Headers:    errorHeaders,
				Body:       internalServerError,
			}
		}

		return nil
	}
}

var badRequest = []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>")
var internalServerError = []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>")
var okRequest = []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>")
