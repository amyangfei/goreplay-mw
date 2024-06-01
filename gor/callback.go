package gor

import "sync"

// Callback is a callback for a given event.
type Callback struct {
	f    func(msg *Message, args map[string]interface{}) *Message
	args map[string]interface{}
}

// ChainCallbacks is a chain of callbacks for a given event.
type ChainCallbacks struct {
	lock sync.Mutex
	m    map[string][]*Callback
}

// DoCallback does the callback for the given event.
func (cc *ChainCallbacks) DoCallback(event string, msg *Message) *Message {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	var resp *Message
	if callbacks, ok := cc.m[event]; ok {
		for _, c := range callbacks {
			r := c.f(msg, c.args)
			if r != nil {
				resp = r
			}
		}
	}
	return resp
}
