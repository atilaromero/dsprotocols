package broadcast

import "github.com/tarcisiocjr/dsprotocols/link"

// BebBroadcastMsg is a message to be broadcasted to all processes.
// Even the sender receives a copy.
type BebBroadcastMsg struct {
	Payload []byte
}

// BebDelivertMsg contains the received brodcast message and the ID of the current process.
type BebDelivertMsg struct {
	ID      int
	Payload []byte
}

// Beb (best effort broadcast) is a struct that contains:
// Numproc: number of known processes.
// Pl: lower level perfect link
// Req: receives beb requests
// Ind: deliver beb messages
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

	go func() {
		for msg, ok := <-beb.Req; ok; msg, ok = <-beb.Req {
			for q := 0; q < numproc; q++ {
				pl.Req <- link.PlSendMsg{
					Dst:     q,
					Payload: msg.Payload,
				}
			}
		}
	}()

	go func() {
		for msg, ok := <-pl.Ind; ok; msg, ok = <-pl.Ind {
			beb.Ind <- BebDelivertMsg{
				ID:      pl.ID,
				Payload: msg.Payload,
			}
		}
	}()

	return beb
}
