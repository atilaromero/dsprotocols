package broadcast

import (
	"github.com/tarcisiocjr/dsprotocols/layer"
	"github.com/tarcisiocjr/dsprotocols/link"
)

type BebBroadcastMsg struct {
	Payload []byte
}
type BebDelivertMsg struct {
	ID      int
	Payload []byte
}

type Beb struct {
	NumProc int
	Pl      link.Pl
	Req     chan BebBroadcastMsg
	Ind     chan BebDelivertMsg
}

func NewBeb(pl link.Pl, numproc int) Beb {
	req := make(chan BebBroadcastMsg)
	ind := make(chan BebDelivertMsg)
	beb := Beb{numproc, pl, req, ind}

	go func(beb Beb) {
		for {
			select {
			case msg := <-beb.Req:
				for q := 0; q < numproc; q++ {
					pl.Req <- link.PlSendMsg{
						Src:     beb.Pl.ID,
						Dst:     q,
						Payload: msg.Payload,
					}
				}
			case msg := <-pl.Ind:
				beb.Ind <- BebDelivertMsg{
					ID:      pl.ID,
					Payload: msg.Payload,
				}
			}
		}
	}(beb)
	return beb
}

type Beb2 struct {
	*layer.Struct
	Pl2 link.Pl2
}

func NewBeb2(ID int, pool map[int]layer.Layer, pl2 link.Pl2) Beb2 {
	var beb Beb2
	beb = Beb2{
		Struct: layer.New(ID, pool),
		Pl2: pl2,
	}
	beb.UpponEvent("BROADCAST", func(msg layer.Msg) {
		for q := 0; q < len(pl2.Pool); q++ {
			pl2.Pool[q].Trigger(ID, "SEND", msg.Payload)
		}
	})
	beb.Pl2.UpponEvent("DELIVER", func(msg layer.Msg){
		beb.Trigger(ID, "DELIVER", msg.Payload)
	})
	return beb
}
