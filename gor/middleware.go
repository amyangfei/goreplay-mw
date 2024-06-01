package gor

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Message represents decoded data from goreplay middleware.
type Message struct {
	ID      string
	Type    string
	Meta    []string
	RawMeta []byte
	HTTP    []byte
}

// Gor is a middleware for goreplay.
type Gor struct {
	callbacks *ChainCallbacks
	stdin     *bufio.Reader
	stdout    *bufio.Writer
}

// NewChainCallbacks returns a new ChainCallbacks.
func NewChainCallbacks() *ChainCallbacks {
	return &ChainCallbacks{
		m: make(map[string][]*Callback),
	}
}

// NewGor returns a new Gor.
func NewGor() *Gor {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	return &Gor{
		callbacks: NewChainCallbacks(),
		stdin:     stdin,
		stdout:    stdout,
	}
}

// On registers a callback for the given event.
func (g *Gor) On(
	event string,
	f func(msg *Message, args map[string]interface{}) *Message,
	idx string,
	args map[string]interface{},
) {
	if idx != "" {
		event += "#" + idx
	}

	g.callbacks.lock.Lock()
	defer g.callbacks.lock.Unlock()

	g.callbacks.m[event] = append(g.callbacks.m[event], &Callback{f: f, args: args})
}

// Run runs the middleware.
func (g *Gor) Run() error {
	for {
		line, err := g.stdin.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				return g.Process(line)
			}
			return nil
		}
		if err = g.Process(line); err != nil {
			return err
		}
	}
}

// Process parses input line, decodes to gor message, trigger necessary
// callbacks and emits result to stdout
func (g *Gor) Process(line string) error {
	msg, err := DecodeGorMsg(line)
	if err != nil {
		return err
	}
	g.Emit(msg)
	return nil
}

// Emit emits the message to stdout.
func (g *Gor) Emit(msg *Message) {
	chanPrefixMap := map[string]string{
		"1": "request",
		"2": "response",
		"3": "replay",
	}
	chanPrefix, ok := chanPrefixMap[msg.Type]
	if !ok {
		return
	}
	resp := msg
	for _, chanID := range []string{"message", chanPrefix, fmt.Sprintf("%s#%s", chanPrefix, msg.ID)} {
		r := g.callbacks.DoCallback(chanID, msg)
		if r != nil {
			resp = r
		}
	}
	g.stdout.WriteString(EncodeGorMsg(resp)) //nolint:errcheck
	g.stdout.WriteString("\n")               //nolint:errcheck
	g.stdout.Flush()
}
