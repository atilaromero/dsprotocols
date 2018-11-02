package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/leadership"
	"github.com/tarcisiocjr/dsprotocols/link"
)

const timeStep = 150 * time.Millisecond

func ExampleNewEc() (map[int]*Ec, map[int]chan leadership.TrustMsg) {
	// first, create a map of perfect links, using process ID as key
	pls := make(map[int]chan<- link.Message)
	pls2 := make(map[int]chan<- link.Message)
	omegas := make(map[int]chan leadership.TrustMsg)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Ecs
	ecs := make(map[int]*Ec)       // each process has a Ec instance, which has a Pl instance, which has a Beb instance
	for i := 0; i < numproc; i++ { // 'i' will be the process ID
		pl := link.NewByChan(i, pls) // when a new Pl is created, it adds itself to pls
		beb := broadcast.NewBeb(pl, numproc)
		pl2 := link.NewByChan(i, pls2)
		omegas[i] = make(chan leadership.TrustMsg)
		ecs[i] = NewEc(pl2, beb, omegas[i], numproc)
		go func(ind chan EcDelivertMsg) {
			for _, ok := <-ind; ok; _, ok = <-ind {
				// fmt.Printf("New Epoch TS: %d, Leader: %v\n", msg.Ts, msg.Leader)
			}
		}(ecs[i].Ind)
	}
	return ecs, omegas
}

func ExampleNewEc_a() {
	ecs, omegas := ExampleNewEc()

	omegas[0] <- leadership.TrustMsg{ID: 0}
	time.Sleep(timeStep)

	for i := 0; i < len(ecs); i++ {
		fmt.Printf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
	}

	// Output:
	// Epoch TS: 3, Leader: 0
	// Epoch TS: 3, Leader: 0
	// Epoch TS: 3, Leader: 0
}

func TestNewEc_b(t *testing.T) {
	ecs, omegas := ExampleNewEc()

	omegas[0] <- leadership.TrustMsg{ID: 0}

	omegas[1] <- leadership.TrustMsg{ID: 1}
	omegas[0] <- leadership.TrustMsg{ID: 1}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	time.Sleep(timeStep)

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Lastts != ecs[0].Lastts {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("TS not unanimous:\n" + s)
			break
		}
	}

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Trusted != 1 {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("Wrong leader:\n" + s)
			break
		}
	}
}

func TestNewEc_c(t *testing.T) {
	ecs, omegas := ExampleNewEc()

	omegas[0] <- leadership.TrustMsg{ID: 0}

	omegas[0] <- leadership.TrustMsg{ID: 1}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	omegas[1] <- leadership.TrustMsg{ID: 1}
	time.Sleep(timeStep)

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Lastts != ecs[0].Lastts {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("TS not unanimous:\n" + s)
			break
		}
	}

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Trusted != 1 {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("Wrong leader:\n" + s)
			break
		}
	}
}

func TestNewEc_d(t *testing.T) {
	ecs, omegas := ExampleNewEc()

	omegas[0] <- leadership.TrustMsg{ID: 0}

	omegas[0] <- leadership.TrustMsg{ID: 2}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	omegas[1] <- leadership.TrustMsg{ID: 0}
	omegas[0] <- leadership.TrustMsg{ID: 1}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	omegas[1] <- leadership.TrustMsg{ID: 2}
	omegas[0] <- leadership.TrustMsg{ID: 2}
	omegas[2] <- leadership.TrustMsg{ID: 0}
	omegas[1] <- leadership.TrustMsg{ID: 0}
	omegas[0] <- leadership.TrustMsg{ID: 1}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	omegas[1] <- leadership.TrustMsg{ID: 2}
	omegas[0] <- leadership.TrustMsg{ID: 2}
	omegas[2] <- leadership.TrustMsg{ID: 0}
	omegas[1] <- leadership.TrustMsg{ID: 0}
	omegas[0] <- leadership.TrustMsg{ID: 1}
	omegas[2] <- leadership.TrustMsg{ID: 1}
	omegas[1] <- leadership.TrustMsg{ID: 0}
	omegas[0] <- leadership.TrustMsg{ID: 0}
	omegas[2] <- leadership.TrustMsg{ID: 0}
	time.Sleep(timeStep)

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Lastts != ecs[0].Lastts {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("TS not unanimous:\n" + s)
			break
		}
	}

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Trusted != 0 {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("Wrong leader:\n" + s)
			break
		}
	}
}

func TestNewEc_e(t *testing.T) {
	ecs, omegas := ExampleNewEc()

	for n := 0; n < 300; n++ {
		for i := 0; i < 3; i++ {
			for l := 0; l < 3; l++ {
				omegas[i] <- leadership.TrustMsg{ID: l}
			}
		}
	}
	time.Sleep(timeStep)

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Lastts != ecs[0].Lastts {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("TS not unanimous:\n" + s)
			break
		}
	}

	for j := 0; j < len(ecs); j++ {
		if ecs[j].Trusted != 2 {
			s := ""
			for i := 0; i < len(ecs); i++ {
				s += fmt.Sprintf("Epoch TS: %d, Leader: %v\n", ecs[i].Lastts, ecs[i].Trusted)
			}
			t.Errorf("Wrong leader:\n" + s)
			break
		}
	}
}
