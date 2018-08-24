package link

import (
	"github.com/tarcisiocjr/dsprotocols/layer"
)

type PlSendMsg struct {
	Src     int
	Dst     int
	Payload []byte
}

type PlDeliverMsg struct {
	Src     int
	Dst     int
	Payload []byte
}

type Pl struct {
	ID  int
	Req chan PlSendMsg
	Ind chan PlDeliverMsg
}

func NewPl(ID int, pls map[int]Pl) Pl {
	req := make(chan PlSendMsg, 1)
	ind := make(chan PlDeliverMsg, 1)
	pl := Pl{ID, req, ind}
	pls[ID] = pl

	go func(pl Pl) {
		for {
			msg := <-pl.Req
			pls[msg.Dst].Ind <- PlDeliverMsg{
				Dst:     msg.Src,
				Src:     msg.Dst,
				Payload: msg.Payload,
			}
		}
	}(pl)
	return pl
}

type Pl2 struct {
	*layer.Struct
}

func NewPl2(ID int, pool map[int]layer.Layer) Pl2 {
	pl := Pl2{
		Struct: layer.New(ID, pool),
	}
	pl.UpponEvent("SEND", func(msg layer.Msg) {
		pool[msg.Dst].Trigger(ID, "DELIVER", msg.Payload)
	})
	return pl
}
