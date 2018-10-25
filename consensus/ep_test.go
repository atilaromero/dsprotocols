package consensus

import (
	"testing"
	"time"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

func TestNewEp(t *testing.T) {

	pls := make(map[int]chan<- link.Message)
	pls2 := make(map[int]chan<- link.Message)
	numproc := 3 // sets the number of known processes
	initialState := State{0, 0}

	// creates and populate a slice of Eps
	eps := make(map[int]*Ep)
	for i := 0; i < numproc; i++ {
		pl := link.NewByChan(i, pls)
		beb := broadcast.NewBeb(pl, numproc)
		pl2 := link.NewByChan(i, pls2)
		eps[i] = NewEp(pl2, beb, numproc, initialState, i == 0)
	}

	for _, ep := range eps {
		ep.Req <- EpProposeMsg{Abort: false, Val: 5}
	}

	//for _, ep := range eps {
	//	ep.Req <- EpProposeMsg{Abort: true, Val: -1}
	//}

	time.Sleep(time.Millisecond * 1000)

}
