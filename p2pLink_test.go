package main

import (
	"fmt"
)

func ExampleNewP2PLink() {
	h := NewCommunicationHandler(func(msg []byte) {
		fmt.Printf("req: %s\n", msg)
	}, func(msg []byte) {
		fmt.Printf("ind: %s\n", msg)
	})

	l := NewP2PLink(h)
	l.Req <- []byte("teste")
	l.Ind <- []byte("teste2")
	close(l.Req)
	close(l.Ind)
	<-l.Closed
	// Output:
	// req: teste
	// ind: teste2
}
