package linkchannel

import "fmt"

func ExampleNewPl() {
	pls := make(map[int]Pl)
	l0 := NewPl(0, pls)
	l1 := NewPl(1, pls)
	l0.Req <- PlSendMsg{
		Dst:     1,
		Payload: []byte("teste"),
	}
	msg := <-l1.Ind
	fmt.Println(msg.Src)
	fmt.Println(string(msg.Payload))
	// Output:
	// 0
	// teste
}
