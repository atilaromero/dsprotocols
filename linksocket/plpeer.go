package linksocket

import (
	"encoding/gob"
	"io"
)

// -----------------------------------------------------------------------------

// Peer is a `struct` representing a remote `PLSocket`.
type Peer struct {
	ID   string
	addr string
	sock *gob.Encoder
}

// NewPeer returns a new `Peer`.
func NewPeer(ID, addr string, fd io.Writer) *Peer {
	return &Peer{ID: ID, addr: addr, sock: gob.NewEncoder(fd)}
}
