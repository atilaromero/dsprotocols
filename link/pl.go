package link

type PlSendMsg struct {
	Src     int
	Dst     int
	Payload []byte
}

type PlDeliverMsg struct {
	Src     int
	Dst     int
	Payload []byte
}

type Pl struct {
	ID  int
	Req chan PlSendMsg
	Ind chan PlDeliverMsg
}

func NewPl(ID int, pls map[int]Pl) Pl {
	req := make(chan PlSendMsg, 1)
	ind := make(chan PlDeliverMsg, 1)
	pl := Pl{ID, req, ind}
	pls[ID] = pl

	go func(pl Pl) {
		for {
			msg := <-pl.Req
			pls[msg.Dst].Ind <- PlDeliverMsg{
				Dst:     msg.Src,
				Src:     msg.Dst,
				Payload: msg.Payload,
			}
		}
	}(pl)
	return pl
}
