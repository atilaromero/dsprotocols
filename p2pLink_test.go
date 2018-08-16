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
	l.req <- []byte("teste")
	l.ind <- []byte("teste2")
	close(l.req)
	close(l.ind)
	<-l.closed
	// Output:
	// req: teste
	// ind: teste2
}
