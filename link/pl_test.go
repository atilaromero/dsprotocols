package link

import "fmt"

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
