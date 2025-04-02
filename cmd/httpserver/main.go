package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	if r.RequestLine.RequestTarget == "/video" {
		handlerGetVideo(w)
		return nil
	}

	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, r)
		return nil
	}

	if r.RequestLine.RequestTarget == "/yourproblem" {
		handler400()
		return nil
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		handler500()
		return nil
	}

	handler200(w)
	return nil
}

func proxyHandler(w response.Writer, r *request.Request) *server.HandlerError {
	target := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")

	res, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", target))
	if err != nil {
		return getUnknownHandlerError(err)
	}
	defer res.Body.Close()

	err = w.WriteStatusLine(response.StatusOk)
	if err != nil {
		getUnknownHandlerError(err)
	}

	resHeaders := headers.NewHeaders()
	resHeaders.Set("Transfer-Encoding", "chunked")
	resHeaders.Set("Connection", "close")
	resHeaders.Set("Trailer", "X-Content-Sha256, X-Content-Length")

	err = w.WriteHeaders(resHeaders)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	buff := make([]byte, 1024)
	fullBody := make([]byte, 0)
	for {
		n, err := res.Body.Read(buff)

		if n > 0 {
			_, err := w.WriteChunkedBody(buff[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}

			fullBody = append(fullBody, buff[:n]...)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}

	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		return getUnknownHandlerError(err)
	}

	trailerHeaders := headers.NewHeaders()
	trailerHeaders.Set("X-Content-Sha256", fmt.Sprintf("%x", sha256.Sum256(fullBody)))
	trailerHeaders.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

	err = w.WriteTrailers(trailerHeaders)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	return nil
}

func handlerGetVideo(w response.Writer) *server.HandlerError {
	err := w.WriteStatusLine(response.StatusOk)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	headers := headers.NewHeaders()
	headers.Set("Content-Type", "video/mp4")

	err = w.WriteHeaders(headers)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	video, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		return getUnknownHandlerError(err)
	}

	_, err = w.WriteBody(video)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	return nil
}

func handler400() *server.HandlerError {
	errorHeaders := headers.NewHeaders()
	errorHeaders.Set("Content-Type", "text/html")

	return &server.HandlerError{
		StatusCode: response.StatusBadRequest,
		Headers:    errorHeaders,
		Body:       badRequest,
	}
}

func handler500() *server.HandlerError {
	errorHeaders := headers.NewHeaders()
	errorHeaders.Set("Content-Type", "text/html")

	return &server.HandlerError{
		StatusCode: response.StatusInternalServerError,
		Headers:    errorHeaders,
		Body:       internalServerError,
	}
}

func handler200(w response.Writer) *server.HandlerError {
	err := w.WriteStatusLine(response.StatusOk)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	headers := response.GetDefaultHeaders(len(okRequest))
	headers.OverrideHeader("Content-Type", "text/html")

	err = w.WriteHeaders(headers)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	_, err = w.WriteBody(okRequest)
	if err != nil {
		return getUnknownHandlerError(err)
	}

	return nil
}

func getUnknownHandlerError(err error) *server.HandlerError {
	errorHeaders := headers.NewHeaders()
	errorHeaders.Set("Content-Type", "text/plain")

	return &server.HandlerError{
		StatusCode: response.StatusInternalServerError,
		Headers:    errorHeaders,
		Body:       []byte(err.Error()),
	}
}

var badRequest = []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>")
var internalServerError = []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>")
var okRequest = []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>")
