package link

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type socketPeer struct {
	addr    string
	conn    *net.Conn
	encoder *gob.Encoder
}

type linkSocket struct {
	id         int
	addr       string
	peers      map[int]socketPeer
	indication chan Message
}

func (pl linkSocket) GetDeliver() <-chan Message {
	return pl.indication
}

func (pl linkSocket) ID() int {
	return pl.id
}

// Send sends a `PLSocket.Message` of type `what` to `pl.peer[to]` at the address
// `addr`. If no connection is reachable at `addr` or if `pl.peer[to]` does not
// exist, the function retries five times and returns an `error` if it does not
// succeed.
func (pl linkSocket) Send(id int, payload []byte) error {
	peer, ok := pl.peers[id]
	if !ok {
		return fmt.Errorf("peer not found: %d", id)
	}
	if peer.conn == nil {
		conn, err := net.Dial("tcp", peer.addr)
		if err != nil {
			return err
		}
		peer.conn = &conn
		peer.encoder = gob.NewEncoder(conn)
	}
	return peer.encoder.Encode(Message{Src: pl.id, Payload: payload})
}

// NewBySocket returns a new `Socket` or an `error`.
// NOTE: All connections to `Peer`s are established during this function.
func NewBySocket(id int, proto string, peers map[int]string) (Link, error) {
	pl := linkSocket{
		id:         id,
		addr:       peers[id],
		peers:      make(map[int]socketPeer),
		indication: make(chan Message),
	}

	for k, addr := range peers {
		pl.peers[k] = socketPeer{
			addr: addr,
		}
	}

	err := pl.listen(proto, peers[id])
	if err != nil {
		return nil, err
	}

	return pl, nil
}

// listen is a helper function that spawns goroutines handling new `Peers`
// connections to `pl`'s socket.
//
// NOTE: this function is an infinite loop.
//
// Listen makes `pl` listens on the address `addr` provided using the protocol
// `proto` and returns an `error` if something occurs.
func (pl *linkSocket) listen(proto, addr string) error {
	ln, err := net.Listen(proto, addr)
	if err != nil {
		return err
	}

	go func(ln net.Listener) {
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func(conn net.Conn) {
				var msg Message
				dec := gob.NewDecoder(conn)
				err = dec.Decode(&msg)
				if err != nil {
					log.Fatal(err)
				}
				pl.indication <- msg
			}(conn)
		}
	}(ln)
	return nil
}
