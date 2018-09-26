package consensus

import (
	"testing"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

func TesteNewEp(t *testing.T) {

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

}
