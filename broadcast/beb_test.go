package broadcast

import (
	"fmt"

	"github.com/tarcisiocjr/dsprotocols/link"
)

func ExampleNewBeb() {
	pls := make(map[int]link.Pl)
	numproc := 3
	bebs := []Beb{}
	for i := 0; i < numproc; i++ {
		pl := link.NewPl(i, pls)
		bebs = append(bebs, NewBeb(pl, numproc))
	}
	bebs[0].Req <- BebBroadcastMsg{
		Payload: []byte("teste"),
	}
	for i := 0; i < numproc; i++ {
		msg := <-bebs[i].Ind
		fmt.Printf("%d: %s\n", msg.ID, string(msg.Payload))
	}
	// Output:
	// 0: teste
	// 1: teste
	// 2: teste
}
