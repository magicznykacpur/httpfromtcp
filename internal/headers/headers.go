package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const crlf = "\r\n"
var specialCharacters = []string{"!", "~", "#", "$", "%", "&", "'", "*", "-", ".", "^", "_", "`", "|"}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	headerText := string(data[:idx])
	parts := strings.Split(strings.TrimSpace(headerText), ": ")

	if len(parts) != 2 || strings.ContainsAny(parts[0], " ") {
		return 0, false, fmt.Errorf("invalid header format")
	}

	for _, r := range parts[0] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if !slices.Contains(specialCharacters, string(r)) {
				return 0, false, fmt.Errorf("invalid header key format")
			}
		}
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	currentValue, ok := h.Get(key)
	if ok {
		h.Set(key, fmt.Sprintf("%s, %s", currentValue, value))
	} else {
		h.Set(key, value)
	}

	return len(data[:idx]) + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[key]
	return val, ok
}
