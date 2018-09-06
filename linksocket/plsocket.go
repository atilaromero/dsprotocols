package linksocket

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// PLSocket our struct per proccess
type PLSocket struct {
	*net.TCPListener

	ID                 string
	addr               string
	leader             string
	PID                int
	peers              Peers
	mu                 *sync.RWMutex
	besteffortChan     chan Message
	reliableChan       chan Message
	leaderChan         chan Message
	epochchangeChan    chan Message
	epochconsensusChan chan Message
	totalorderChan     chan Message
}

// NewSocket returns a new `Socket` or an `error`.
// NOTE: All connections to `Peer`s are established during this function.
func NewSocket(ID, addr, proto string, pid int, peers map[string]string) (*PLSocket, error) {
	pl := &PLSocket{
		ID:                 ID,
		addr:               addr,
		leader:             ID,
		PID:                pid,
		peers:              NewPeerMap(),
		mu:                 &sync.RWMutex{},
		besteffortChan:     make(chan Message, 0),
		reliableChan:       make(chan Message, 1),
		leaderChan:         make(chan Message, 0),
		epochchangeChan:    make(chan Message, 0),
		epochconsensusChan: make(chan Message, 0),
		totalorderChan:     make(chan Message, 0),
	}

	if err := pl.Listen(proto, addr); err != nil {
		return nil, err
	}

	if err := pl.Connect(proto, peers); err != nil {
		return nil, err
	}
	return pl, nil
}

func (pl *PLSocket) NewBebsocket() {
	teste := []byte{'R', 'e', 'n', 'a', 'n'}
	for _, rpl := range pl.peers.PeerData() {
		_ = pl.Send(rpl.ID, rpl.Addr, 0, teste)
	}
}

func (pl *PLSocket) receive(rwc io.ReadWriteCloser) {
	var msg Message
	dec := gob.NewDecoder(rwc)

	for {
		if err := dec.Decode(&msg); err == io.EOF || msg.Type == CLOSE {
			_ = rwc.Close()
			//if msg.PeerID > 0 {
			//if msg.PeerID == pl.Coordinator() {
			//	pl.peers.Delete(msg.PeerID)
			//	pl.SetCoordinator(pl.ID)
			//	pl.Elect()
			//}
			break
		} else if msg.Type == 10000 {
			select {
			case pl.leaderChan <- msg:
				continue
			case <-time.After(10 * time.Millisecond):
				continue
			}
		} else {
			pl.besteffortChan <- msg
			fmt.Println("recebi a informacao")
		}
	}
}

// listen is a helper function that spawns goroutines handling new `Peers`
// connections to `pl`'s socket.
//
// NOTE: this function is an infinite loop.
func (pl *PLSocket) listen() {
	for {
		conn, err := pl.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go pl.receive(conn)
	}
}

// Listen makes `pl` listens on the address `addr` provided using the protocol
// `proto` and returns an `error` if something occurs.
func (pl *PLSocket) Listen(proto, addr string) error {
	laddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return err
	}
	pl.TCPListener, err = net.ListenTCP(proto, laddr)
	if err != nil {
		return err
	}
	go pl.listen()
	return nil
}

// connect is a helper function that resolves the tcp address `addr` and try
// to establish a tcp connection using the protocol `proto`. The established
// connection is set to `pl.peers[ID]` or the function returns an `error`
// if something occurs.
func (pl *PLSocket) connect(proto, addr, ID string) error {
	raddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return err
	}
	sock, err := net.DialTCP(proto, nil, raddr)
	if err != nil {
		return err
	}

	pl.peers.Add(ID, addr, sock)
	return nil
}

// Connect performs a connection to the remote `Peer`s and returns an `error`
// if something occurs during connection.
func (pl *PLSocket) Connect(proto string, peers map[string]string) (err error) {
	for ID, addr := range peers {
		if pl.ID == ID {
			continue
		}
		if err := pl.connect(proto, addr, ID); err != nil {
			continue
		}
	}
	return nil
}

// Send sends a `PLSocket.Message` of type `what` to `pl.peer[to]` at the address
// `addr`. If no connection is reachable at `addr` or if `pl.peer[to]` does not
// exist, the function retries five times and returns an `error` if it does not
// succeed.
func (pl *PLSocket) Send(to, addr string, what int, payload []byte) error {
	maxRetries := 5

	if !pl.peers.Find(to) {
		_ = pl.connect("tcp4", addr, to)
	}

	for attempts := 0; ; attempts++ {
		err := pl.peers.Write(to, &Message{PeerID: pl.ID, Addr: pl.addr, Type: what, Payload: payload})
		if err == nil {
			break
		}
		if attempts > maxRetries && err != nil {
			return err
		}
		_ = pl.connect("tcp4", addr, to)
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

// NewBebsocket function
