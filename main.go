package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
	"github.com/tarcisiocjr/dsprotocols/linksocket"
)

const conf = "conf.json"

type configuration struct {
	Verbose bool
	Hosts   []string
}

func main() {
	var operation bool
	var allClients map[*Client]int
	allClients := make(map[*linksocket.Client]ints)
	hosts := make(map[int]Host)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [socket|channel]\n", os.Args[0])
		os.Exit(1)
	}

	if strings.Compare(os.Args[1], "channel") == 0 {
		operation = false
	} else {
		operation = true
	}

	if operation == true {
		file, fileErr := os.Open(conf)
		if fileErr != nil {
			getOut(fileErr)
		}

		decoder := json.NewDecoder(file)
		conf := configuration{}
		err := decoder.Decode(&conf)
		if err != nil {
			fmt.Println("Error:", err)
		}

		for ix, ipPort := range conf.Hosts {
			var arr [20]byte
			tmp := strings.Split(ipPort, ":")
			copy(arr[:], tmp[0])
			port, _ := strconv.Atoi(tmp[1])
			pid := os.Getpid()
			host := host{ix, arr, port, pid}
			hosts[ix] = host
			//fmt.Println(ix, tmp[0], port, arr)
			linksocket.NewPlSocket(allClients, host)
		}

		for _, value := range hosts {
			fmt.Printf("%d %v %d %d\n", value.ID, value.IP, value.port, value.pid)
		}
	} else {
		// first, create a map of perfect links, using process ID as keys
		pls := make(map[int]link.Pl)
		numproc := 3 // sets the number of known processes

		// creates and populate a slice of Bebs
		bebs := []broadcast.Beb{}      // each process has a Beb instance, which has a Pl instance
		for i := 0; i < numproc; i++ { // 'i' will be the process ID
			pl := link.NewPl(i, pls) // when a new Pl is created, it adds itself to pls
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

func getOut(fileErr error) {
	fmt.Printf("Cant open %s", conf)
	os.Exit(1)
}
