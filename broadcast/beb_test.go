package broadcast

import (
	"fmt"
	"testing"

	"github.com/tarcisiocjr/dsprotocols/link"
)

func ExampleNewBeb() {
	// first, create a map of perfect links, using process ID as key
	pls := make(map[int]chan<- link.Message)
	numproc := 3 // sets the number of known processes

	// creates and populate a slice of Bebs
	bebs := []Beb{}                // each process has a Beb instance, which has a Pl instance
	for i := 0; i < numproc; i++ { // 'i' will be the process ID
		pl := link.NewByChan(i, pls) // when a new Pl is created, it adds itself to pls
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

func TestNewBebWithSocket(t *testing.T) {
	// first, create a map of perfect links, using process ID as key
	pls := make(map[int]string)
	numproc := 3 // sets the number of known processes

	for i := 0; i < numproc; i++ { // 'i' will be the process ID
		pls[i] = fmt.Sprintf(":%d", 10000+i)
	}

	// creates and populate a slice of Bebs
	bebs := []Beb{}                // each process has a Beb instance, which has a Pl instance
	for i := 0; i < numproc; i++ { // 'i' will be the process ID
		pl, err := link.NewBySocket(i, "tcp4", pls) // when a new Pl is created, it adds itself to pls
		if err != nil {
			t.Errorf("error creating link: %v", err)
		}
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
		if msg.ID != i {
			t.Errorf("ID: expected %d; got %d", i, msg.ID)
		}
		if string(msg.Payload) != "test" {
			t.Errorf("payload: expected %s; got %s", "test", string(msg.Payload))
		}
	}
}
