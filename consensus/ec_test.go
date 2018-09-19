package consensus

import (
	"fmt"
	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

func ExampleNewEc() {
	// first, create a map of perfect links, using process ID as key
	pls := make(map[int]chan<- link.Message)
	pls2 := make(map[int]chan<- link.Message)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Ecs
	ecs := []Ec{}                  		// each process has a Ec instance, which has a Pl instance, which has a Beb instance
	for i := 0; i < numproc; i++ { 		// 'i' will be the process ID
		pl := link.NewByChan(i, pls) 	// when a new Pl is created, it adds itself to pls
		beb := broadcast.NewBeb(pl, numproc)
		pl2 := link.NewByChan(i, pls2)
		ecs = append(ecs, NewEc(pl2, beb, nil, 0))
	}

	// propose a epoch-change
	newts1 := 100
	// ask a process to broadcast a message
	ecs[0].Beb.Req <- broadcast.BebBroadcastMsg{
		Payload: []byte{byte(newts1)},
	}

	// wait until all processes have treated the broadcast message
	for i := 0; i < numproc; i++ {
		msg := <-ecs[i].Ind // the resulting broadcasted message is read from indication channel
		// our treatment is just to print the message
		fmt.Printf("New Epoch TS: %d, Leader: %v\n", msg.Ts, msg.Leader)
	}
	// Output:
	// New Epoch TS: 100, Leader: 0
	// New Epoch TS: 100, Leader: 0
	// New Epoch TS: 100, Leader: 0
}
