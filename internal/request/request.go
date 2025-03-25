package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, fmt.Errorf("couldn't read bytes: %v", err)
	}

	parts := strings.Split(string(bytes), "\r\n")
	methodParts := strings.Split(parts[0], " ")

	if len(methodParts) != 3 {
		return &Request{}, fmt.Errorf("invalid method request line")
	}

	method := methodParts[0]
	requestTarget := methodParts[1]
	protocolVersion := methodParts[2]

	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return &Request{}, fmt.Errorf("invalid http method")
		}
	}

	if protocolVersion != "HTTP/1.1" {
		return &Request{}, fmt.Errorf("invalid protocol version")
	}

	methodRequestLine := RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   strings.Split(protocolVersion, "/")[1],
	}

	return &Request{RequestLine: methodRequestLine}, nil
}
