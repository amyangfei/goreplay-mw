package gor

import (
	"encoding/hex"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Counter struct {
	Value int
}

func (c *Counter) Incr() {
	c.Value++
}

func incrRecv(msg *Message, args map[string]interface{}) *Message {
	c := args["passby"].(*Counter)
	c.Incr()
	return nil
}

func TestRun(t *testing.T) {
	payload := strings.Join([]string{
		hex.EncodeToString([]byte("1 2 3\nGET / HTTP/1.1\r\n\r\n")),
		hex.EncodeToString([]byte("2 2 3\nHTTP/1.1 200 OK\r\n\r\n")),
		hex.EncodeToString([]byte("2 3 3\nHTTP/1.1 200 OK\r\n\r\n")),
	}, "\n")
	testWithStdCapture(t, payload, func() {
		g := NewGor()
		counter := &Counter{Value: 100}
		args := map[string]interface{}{
			"passby": counter,
		}
		g.On("message", incrRecv, "", args)
		g.On("request", incrRecv, "", args)
		g.On("response", incrRecv, "2", args)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := g.Run()
			require.Nil(t, err)
		}()

		wg.Wait()
		require.Equal(t, 105, counter.Value)
	})
}

func testWithStdCapture(t *testing.T, input string, testFunc func()) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(input))
	require.Nil(t, err)
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	testFunc()
}
