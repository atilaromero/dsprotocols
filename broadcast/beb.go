package broadcast

import "github.com/tarcisiocjr/dsprotocols/linkchannel"

// BebBroadcastMsg is a message to be broadcasted to all processes.
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
	Pl      linkchannel.Pl
	Req     chan BebBroadcastMsg
	Ind     chan BebDelivertMsg
}

// NewBeb returns a Beb struct, which implements Best Effort Broadcast.
// The sender also receives a copy of broadcasted messages.
//
// There are 4 channels here: the 2 channels from the perfect link
// are used for inter process communication, while req and ind are used to
// start and finish the broadcast.
//
// New broadcasts are initiated sending a message to the req channel.
//   Ex: beb.Req <- msg
//
// When an ongoing broadcast is received from another process,
// a deliver is triggered through the ind channel.
//   Ex: msg <- beb.Ind
//
// When using Beb, remember to create a go routine reading from the end channel beb.Ind and
// treating incomming messages.
func NewBeb(pl linkchannel.Pl, numproc int) Beb {
	req := make(chan BebBroadcastMsg)
	ind := make(chan BebDelivertMsg)
	beb := Beb{numproc, pl, req, ind}

	go func() {
		// on new broadcast request
		for msg, ok := <-beb.Req; ok; msg, ok = <-beb.Req {
			// send the message to all known processes through each one's perfect link
			for q := 0; q < numproc; q++ {
				pl.Req <- linkchannel.PlSendMsg{
					Dst:     q,
					Payload: msg.Payload,
				}
			}
		}
	}()

	go func() {
		// when receiving a broadcast from another process
		for msg, ok := <-pl.Ind; ok; msg, ok = <-pl.Ind {
			// deliver the message one layer up
			beb.Ind <- BebDelivertMsg{
				ID:      pl.ID,
				Payload: msg.Payload,
			}
		}
	}()

	return beb
}
