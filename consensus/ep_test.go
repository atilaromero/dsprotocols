package consensus

import (
	"fmt"
	"testing"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

func TestNewEp(t *testing.T) {

	pls := make(map[int]chan<- link.Message)
	pls2 := make(map[int]chan<- link.Message)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Eps
	eps := make(map[int]*Ep)
	for i := 0; i < numproc; i++ {
		pl := link.NewByChan(i, pls)
		beb := broadcast.NewBeb(pl, numproc)
		pl2 := link.NewByChan(i, pls2)
		eps[i] = NewEp(pl2, beb, numproc)
	}

	ets := 0
	leader := 1
	for i, ep := range eps {
		ep.Init(ets, leader)
		ep.Req <- EpProposeMsg{Abort: false, Val: i + 100}
	}

	// time.Sleep(time.Millisecond * 30)
	for _, ep := range eps {
		msg := <-ep.Ind
		expect := 101
		fmt.Println("Result: ", msg.Val)
		if msg.Val != expect {
			t.Errorf("Wrong result. Expected %d, got %d", expect, msg.Val)
		}
	}

	for _, ep := range eps {
		ep.Req <- EpProposeMsg{Abort: true, Val: 200}
		<-ep.Ind
	}
	// time.Sleep(time.Millisecond * 100)

	ets = 1
	leader = 2
	for i, ep := range eps {
		ep.Init(ets, leader)
		ep.Req <- EpProposeMsg{Abort: false, Val: 1000 + i}
	}

	// time.Sleep(time.Millisecond * 300)
	for _, ep := range eps {
		msg := <-ep.Ind
		expect := 1002
		fmt.Println("Result: ", msg.Val)
		if msg.Val != expect {
			t.Errorf("Wrong result. Expected %d, got %d", expect, msg.Val)
		}
	}
}
