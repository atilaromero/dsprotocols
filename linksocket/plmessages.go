package linksocket

// -----------------------------------------------------------------------------

// Message Type.
const (
	ELECTION = iota
	OK
	LEADER
	CLOSE
	EC
	EP
	UTOB
)

// -----------------------------------------------------------------------------

// Message is a `struct` used for communication between `process`s.
type Message struct {
	PeerID  string
	Addr    string
	Type    int
	Payload []byte
}
