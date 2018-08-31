package broadcast

import (
	"fmt"

	"github.com/tarcisiocjr/dsprotocols/link"
)

func ExampleNewBeb() {
	// first, create a map of perfect links, using process ID as key
	pls := make(map[int]link.Pl)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Bebs
	bebs := []Beb{}                // each process has a Beb instance, which has a Pl instance
	for i := 0; i < numproc; i++ { // 'i' will be the process ID
		pl := link.NewPl(i, pls) // when a new Pl is created, it adds itself to pls
		bebs = append(bebs, NewBeb(pl, numproc))
	}

	// ask a process to broadcast a message to the others
	bebs[0].Req <- BebBroadcastMsg{
		Payload: []byte("test"),
	}

	// wait until all processes have treated the broadcast message
	for i := 0; i < numproc; i++ {
		msg := <-bebs[i].Ind // the resulting broadcasted message is read from indication channel

		// our treatment is just to print the message
		fmt.Printf("%d: %s\n", msg.ID, string(msg.Payload))
	}
	// Output:
	// 0: test
	// 1: test
	// 2: test
}
