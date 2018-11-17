package consensus

import (
	"fmt"
	"testing"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/leadership"
	"github.com/tarcisiocjr/dsprotocols/link"
)

func TestNewUc(t *testing.T) {

	plsEcs := make(map[int]chan<- link.Message)
	pls2Ecs := make(map[int]chan<- link.Message)
	plsEps := make(map[int]chan<- link.Message)
	pls2Eps := make(map[int]chan<- link.Message)
	omegas := make(map[int]chan leadership.TrustMsg)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Ucs
	ucs := make(map[int]*Uc)
	for i := 0; i < numproc; i++ {

		plEc := link.NewByChan(i, plsEcs) // when a new Pl is created, it adds itself to pls
		bebEc := broadcast.NewBeb(plEc, numproc)
		pl2Ec := link.NewByChan(i, pls2Ecs)
		omegas[i] = make(chan leadership.TrustMsg)
		ec := NewEc(pl2Ec, bebEc, omegas[i], numproc)

		plEp := link.NewByChan(i, plsEps)
		bebEp := broadcast.NewBeb(plEp, numproc)
		pl2Ep := link.NewByChan(i, pls2Eps)
		ep := NewEp(pl2Ep, bebEp, numproc)

		ucs[i] = NewUC(ec, ep)
		ucs[i].Init()
	}

	for i, uc := range ucs {
		uc.Req <- UcProposeMsg{Val: i + 100}
	}

	for _, uc := range ucs {
		msg := <-uc.Ind
		expect := 100
		fmt.Println("Result: ", msg.Val)
		if msg.Val != expect {
			t.Errorf("Wrong result. Expected %d, got %d", expect, msg.Val)
		}
	}

}
