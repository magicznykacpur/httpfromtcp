package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// valid single header with extra whitespaces
	headers = NewHeaders()
	data = []byte("      Host: localhost:42069            \r\n\r\n")

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 41, n)
	assert.False(t, done)

	// valid 2 headers with existing headers
	headers = NewHeaders()
	host := []byte("Host: localhost:42069\r\n\r\n")
	authorization := []byte("Authorization: Bearer token\r\n\r\n")

	n, done, err = headers.Parse(host)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(authorization)
	require.NoError(t, err)
	assert.Equal(t, "Bearer token", headers["authorization"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// valid done
	headers = NewHeaders()
	data = []byte("\r\n\r\n")

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.True(t, done)
	assert.Equal(t, 2, n)

	// invalid spacing header
	headers = NewHeaders()
	data = []byte("         Host  : localhost:42069            \r\n\r\n")

	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, "invalid header format", err.Error())
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// invalid special key character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, "invalid header key format", err.Error())
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// valid special key character
	headers = NewHeaders()
	data = []byte("H!~st: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["h!~st"])
	assert.Equal(t, 24, n)
	assert.False(t, done)

	// multiple values for one key
	headers = NewHeaders()
	data = []byte("Set-Person: kacpi\r\n\r\n")
	moreData := []byte("Set-Person: mati\r\n\r\n")
	evenMoreData := []byte("Set-Person: slawek\r\n\r\n")

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "kacpi", headers["set-person"])
	assert.Equal(t, 19, n)
	assert.False(t, done)
	
	n, done, err = headers.Parse(moreData)
	require.NoError(t, err)
	assert.Equal(t, "kacpi, mati", headers["set-person"])
	assert.Equal(t, 18, n)
	assert.False(t, done)
	
	n, done, err = headers.Parse(evenMoreData)
	require.NoError(t, err)
	assert.Equal(t, "kacpi, mati, slawek", headers["set-person"])
	assert.Equal(t, 20, n)
	assert.False(t, done)
}
