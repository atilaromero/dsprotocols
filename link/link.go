package link

// Message is use for indication.
// These messages are intended to be read by the upper layer.
// Src is the sender's process ID in the link pool.
type Message struct {
	Src     int
	Payload []byte
}

// Link is a struct that contains the ID of a process and 2 channels:
// Link receives messages from a upper layer through Req
// and deliver messages through Ind.
type Link interface {
	Send(id int, payload []byte) error
	GetDeliver() <-chan Message
	ID() int
}
