package gor

import (
	"bytes"
	"encoding/hex"
	"strings"
)

// EncodeGorMsg encodes a message into a string that can be sent to gor.
func EncodeGorMsg(msg *Message) string {
	hexEncode := func(data []byte) []byte {
		encoded := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(encoded, data)
		return encoded
	}

	var buf bytes.Buffer
	buf.Write(hexEncode(msg.RawMeta))
	buf.Write(hexEncode([]byte("\n")))
	buf.Write(hexEncode(msg.HTTP))
	return buf.String()
}

// DecodeGorMsg decodes a message from a string that was sent by gor.
func DecodeGorMsg(line string) (*Message, error) {
	line = strings.TrimSpace(line)
	payload, err := hex.DecodeString(line)
	if err != nil {
		return nil, err
	}
	metaEnd := bytes.IndexByte(payload, '\n')
	meta := strings.Split(string(payload[:metaEnd]), " ")
	return &Message{
		ID:      meta[1],
		Type:    meta[0],
		Meta:    meta,
		RawMeta: payload[:metaEnd],
		HTTP:    payload[metaEnd+1:],
	}, nil
}
