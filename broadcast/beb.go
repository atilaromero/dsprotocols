package broadcast

import "github.com/tarcisiocjr/dsprotocols/link"

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
						Src: beb.Pl.ID,
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
