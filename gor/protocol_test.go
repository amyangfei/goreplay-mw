package gor

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	payload := hex.EncodeToString([]byte("1 2 3\nGET / HTTP/1.1\r\n\r\n"))
	expected := &Message{
		Type:    "1",
		ID:      "2",
		Meta:    []string{"1", "2", "3"},
		RawMeta: []byte("1 2 3"),
		HTTP:    []byte("GET / HTTP/1.1\r\n\r\n"),
	}
	msg, err := DecodeGorMsg(payload)
	require.Nil(t, err)
	require.Equal(t, expected, msg)

	data := EncodeGorMsg(msg)
	require.Equal(t, payload, data)
}
