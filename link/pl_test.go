package link

import (
	"time"
	"github.com/tarcisiocjr/dsprotocols/layer"
	"fmt"
)

func ExampleNewPl() {
	pls := make(map[int]Pl)
	l0 := NewPl(0, pls)
	l1 := NewPl(1, pls)
	l0.Req <- PlSendMsg{
		Src:     0,
		Dst:     1,
		Payload: []byte("teste"),
	}
	msg := <-l1.Ind
	fmt.Println(string(msg.Payload))
	// Output:
	// teste
}


func ExampleNewPl2() {
	pls := make(map[int]layer.Layer)
	// l0 := NewPl2(0, pls)
	l1 := NewPl2(1, pls)
	l1.Trigger(0, "SEND", []byte("payload"))
	l1.UpponEvent("DELIVER", func(msg layer.Msg){
		fmt.Println(string(msg.Payload))
	})
	time.Sleep(time.Second)
	// Output:
	// payload
}

