package consensus

import (
	"fmt"
	"testing"
	"time"

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

func CreateUCs(numproc int) map[int]*Uc {
	plsEcs := make(map[int]chan<- link.Message)
	pls2Ecs := make(map[int]chan<- link.Message)
	plsEps := make(map[int]chan<- link.Message)
	pls2Eps := make(map[int]chan<- link.Message)
	omegas := make(map[int]chan leadership.TrustMsg)

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
	return ucs
}

func TestNewUc_b(t *testing.T) {
	numproc := 3 // sets the number of known processes
	for instance := 0; instance < 100; instance++ {
		ucs := CreateUCs(numproc)
		for i, uc := range ucs {
			uc.Req <- UcProposeMsg{Val: i + instance}
		}

		for _, uc := range ucs {
			msg := <-uc.Ind
			expect := instance
			fmt.Println("Result: ", msg.Val)
			if msg.Val != expect {
				t.Errorf("Wrong result. Expected %d, got %d", expect, msg.Val)
			}
		}

		for _, uc := range ucs {
			// close(uc.EcInstance.LeaderDetector)
			close(uc.EpInstance.Req)
			close(uc.EpInstance.Beb.Req)
			close(uc.EcInstance.Beb.Req)
		}
	}
	time.Sleep(200 * time.Millisecond)
}
