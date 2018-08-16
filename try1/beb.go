package try1

import (
	"reflect"
)

type PerfectLink struct {
	req, ind chan struct{}
}

type PerfectP2PLink struct {
	pl []PerfectLink
}

func (p2p *PerfectP2PLink) send(process int, msg struct{}) {
	p2p.pl[process].req <- msg
}

func (p2p *PerfectP2PLink) receive() (int, []byte, bool) {
	// dynamic select case, from
	// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
	cases := make([]reflect.SelectCase, len(p2p.pl))
	for i, pl := range p2p.pl {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(pl.ind)}
	}
	chosen, msg, ok := reflect.Select(cases)
	return chosen, msg.Bytes(), ok // []byte or struct{}? How do we convert msg to struct{}?
}

type Beb struct {
	req, ind chan struct{}
	pl       []PerfectLink
}

func NewBeb(pl []PerfectLink) Beb {
	req := make(chan struct{})
	ind := make(chan struct{})
	go func() {
		for {
		}
	}()

	return Beb{
		req: req,
		ind: ind,
		pl:  pl,
	}
}

func (b *Beb) send(msg struct{}) {
	b.req <- msg
}

func (b *Beb) receive() struct{} {
	return <-b.ind
}
