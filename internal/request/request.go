package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8
const stateInitalized = 0
const stateDone = 1
const crlf = "\r\n"

type Request struct {
	RequestLine RequestLine
	parserState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.parserState {
	case stateInitalized:
		requestLine, parsedBytes, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if parsedBytes == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.parserState = stateDone

		return 0, nil
	case stateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := Request{parserState: stateInitalized}

	for request.parserState != stateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numReadbytes, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.parserState = stateDone
				break
			}

			return nil, err
		}
		readToIndex += numReadbytes

		numParsedBytes, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numParsedBytes:])
		readToIndex -= numParsedBytes
	}

	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	methodParts := strings.Split(str, " ")
	if len(methodParts) != 3 {
		return nil, fmt.Errorf("invalid method request line")
	}

	method := methodParts[0]
	requestTarget := methodParts[1]
	protocolVersion := methodParts[2]

	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, fmt.Errorf("invalid http method")
		}
	}

	if protocolVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid protocol version")
	}

	methodRequestLine := RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   strings.Split(protocolVersion, "/")[1],
	}

	return &methodRequestLine, nil
}
