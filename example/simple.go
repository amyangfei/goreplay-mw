package main

import (
	"fmt"

	"github.com/amyangfei/goreplay-mw/gor"
)

func onRequest(msg *gor.Message, args map[string]interface{}) *gor.Message {
	fmt.Printf("RawMeta: %s\n", msg.RawMeta)
	fmt.Printf("args: %v\n", args)
	return nil
}

func main() {
	g := gor.NewGor()
	g.On("request", onRequest, "", nil)
}
