package broadcast

import (
	"time"
	"github.com/tarcisiocjr/dsprotocols/layer"
	"github.com/tarcisiocjr/dsprotocols/link"
	"fmt"
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

func ExampleNewBeb2() {
	plPool := make(map[int]layer.Layer)
	bebPool := make(map[int]layer.Layer)
	numproc := 3
	for i := 0; i < numproc; i++ {
		pl :=link.NewPl2(i, plPool)
		beb := NewBeb2(i, bebPool, pl)
		beb.UpponEvent("DELIVER", func(msg layer.Msg){
			fmt.Printf("%s\n", string(msg.Payload))
		})
		plPool[i] = pl
		bebPool[i] = beb
	}
	bebPool[0].Trigger(0, "BROADCAST", []byte("teste"))
	time.Sleep(time.Second)
	// Output:
	// teste
	// teste
	// teste
}
