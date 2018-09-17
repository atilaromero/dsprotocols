package link

type linkChannel struct {
	id         int
	indication chan Message
	peers      map[int]chan<- Message
}

func (l linkChannel) Send(id int, payload []byte) error {
	// find destination in pool and fill Src with current process ID
	go func() {
		l.peers[id] <- Message{
			Src:     l.id,
			Payload: payload,
		}
	}()
	return nil
}

func (l linkChannel) GetDeliver() <-chan Message {
	return l.indication
}

func (l linkChannel) ID() int {
	return l.id
}

// NewPl returns a new Pl struct. This implementations uses
// only channels, instead of sockets, designed to be used in a non distributed environment (single machine).
//
// For a network implementation with similar behavior, see XXX...XXX
//
// ID refers to the process than owns this Pl.
// peers is a, possibly empty, pool of all Pl. NewPl will modify this pool, incluing the new Pl in it.
//
// The requisition channel is used to receive messages from the upper layer.
//   Ex: pl.Req <- msg
// When a new message is received through Req,
// that message's destination is seeked in the pool of perfect links (peers), and sent through the
// indication channel.
// Is up to the upper layer to collect messages from the indication channel.
//   Ex: msg <- pl.Ind
//
// The indication channel will block if those messages aren't treated somehow.
// To prevent that, create a go routine that continually reads from pl.Ind
func NewByChan(id int, peers map[int]chan<- Message) Link {
	ind := make(chan Message)
	l := linkChannel{id, ind, peers}

	// adds the new created Pl instance to the pool of peers
	peers[id] = ind

	return l
}
