package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/linkchannel"
	"github.com/tarcisiocjr/dsprotocols/linksocket"
)

var confpeerAddr map[string]string

func main() {

	confPeerAddr := map[string]string{
		"1": "0.0.0.0:9991",
		"2": "0.0.0.0:9992",
		"3": "0.0.0.0:9993",
		"4": "0.0.0.0:9994",
	}

	if len(os.Args) < 2 {
		log.Fatal("[Err] node ID required such as '1'")
	}

	pl, err := linksocket.NewSocket(os.Args[1], confPeerAddr[os.Args[1]], "tcp4", os.Getpid(), confPeerAddr)

	if pl == nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	} else {
		// first, create a map of perfect links, using process ID as keys
		pls := make(map[int]linkchannel.Pl)
		numproc := 3 // sets the number of known processes

		// creates and populate a slice of Bebs
		bebs := []broadcast.Beb{}      // each process has a Beb instance, which has a Pl instance
		for i := 0; i < numproc; i++ { // 'i' will be the process ID
			pl := linkchannel.NewPl(i, pls) // when a new Pl is created, it adds itself to pls
			bebs = append(bebs, broadcast.NewBeb(pl, numproc))
		}

		// ask a process to broadcast a message to the others
		bebs[0].Req <- broadcast.BebBroadcastMsg{
			Payload: []byte("test"),
		}

		// wait until all processes have treated the broadcast message
		for i := 0; i < numproc; i++ {
			msg := <-bebs[i].Ind // the resulting broadcasted message is read from indication channel

			// our treatment is just to print the message
			fmt.Printf("%d: %s\n", msg.ID, string(msg.Payload))
		}
	}
}
