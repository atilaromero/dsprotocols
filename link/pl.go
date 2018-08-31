package link

// PlSendMsg is used for requests.
// These messages are intended to be created by the upper layer and sent to Pl.
// Dst is the destination process ID in the link pool.
type PlSendMsg struct {
	Dst     int
	Payload []byte
}

// PlDeliverMsg is use for indication.
// These messages are intended to be read by the upper layer.
// Src is the sender's process ID in the link pool.
type PlDeliverMsg struct {
	Src     int
	Payload []byte
}

// Pl is a struct that contains the ID of a process and 2 channels:
// Pl receives messages from a upper layer through Req
// and deliver messages through Ind.
type Pl struct {
	ID  int
	Req chan PlSendMsg
	Ind chan PlDeliverMsg
}

// NewPl returns a new Pl struct. This implementations uses
// only channels, instead of sockets, designed to be used in a non distributed environment (single machine).
//
// For a network implementation with similar behavior, see XXX...XXX
//
// ID refers to the process than owns this Pl.
// pls is a, possibly empty, pool of all Pl. NewPl will modify this pool, incluing the new Pl in it.
//
// The requisition channel is used to receive messages from the upper layer.
//   Ex: pl.Req <- msg
// When a new message is received through Req,
// that message's destination is seeked in the pool of perfect links (pls), and sent through the
// indication channel.
// Is up to the upper layer to collect messages from the indication channel.
//   Ex: msg <- pl.Ind
//
// The indication channel will block if those messages aren't treated somehow.
// To prevent that, create a go routine that continually reads from pl.Ind
func NewPl(ID int, pls map[int]Pl) Pl {
	req := make(chan PlSendMsg)
	ind := make(chan PlDeliverMsg)
	pl := Pl{ID, req, ind}

	// adds the new created Pl instance to the pool of Pls
	pls[ID] = pl

	go func() {
		// monitor req channel for new messages and deliver them, until channel is closed
		for msg, ok := <-req; ok; msg, ok = <-req {
			// find destination in pool and fill Src with current process ID
			pls[msg.Dst].Ind <- PlDeliverMsg{
				Src:     ID,
				Payload: msg.Payload,
			}
		}
	}()
	return pl
}
