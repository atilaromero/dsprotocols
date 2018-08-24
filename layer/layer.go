package layer

import (
	"fmt"
)

type Msg struct {
	Src     int
	Dst     int
	Type    string
	Payload []byte
}

type Struct struct {
	ID       int
	Events   chan Msg
	Pool     map[int]Layer
	handlers map[string]func(Msg)
}

func New(id int, pool map[int]Layer) *Struct {
	events := make(chan Msg, 1)
	proc := &Struct{id, events, pool, map[string]func(Msg){}}
	pool[id] = proc

	go func() {
		for msg, ok := <-events; ok; msg, ok = <-events {
			h, ok := proc.handlers[msg.Type]
			if !ok {
				fmt.Printf("msg type not known: %s", msg.Type)
			} 
			h(msg)
		}
	}()
	return proc
}

type Layer interface {
	UpponEvent(ev string, f func(Msg))
	Trigger(src int, ev string, payload []byte)
}

func (p *Struct) UpponEvent(ev string, f func(Msg)) {
	p.handlers[ev] = f
}

func (p *Struct) Trigger(src int, ev string, payload []byte) {
	p.Events <- Msg{
		Src:     src,
		Dst:     p.ID,
		Type:    ev,
		Payload: payload,
	}
}
