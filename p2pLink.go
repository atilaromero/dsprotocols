package main

type P2PLink struct {
	Req, Ind chan []byte
	Closed   chan bool
}

func NewP2PLink(h CommunicationHandler) P2PLink {
	req := make(chan []byte)
	ind := make(chan []byte)
	closed := make(chan bool)
	go func() {
		defer close(closed)
		for {
			if req == nil && ind == nil {
				return
			}
			select {
			case msg, ok := <-req:
				if !ok {
					req = nil
					continue
				}
				h.ReqHandler(msg)
			case msg, ok := <-ind:
				if !ok {
					ind = nil
					continue
				}
				h.IndHandler(msg)
			}
		}
	}()
	return P2PLink{
		Req:    req,
		Ind:    ind,
		Closed: closed,
	}
}

type CommunicationHandler interface {
	ReqHandler([]byte)
	IndHandler([]byte)
}
type handler struct {
	r func([]byte)
	i func([]byte)
}

func (h handler) ReqHandler(msg []byte) {
	h.r(msg)
}
func (h handler) IndHandler(msg []byte) {
	h.i(msg)
}
func NewCommunicationHandler(reqHandler func([]byte), indHandler func([]byte)) CommunicationHandler {
	return handler{reqHandler, indHandler}
}
