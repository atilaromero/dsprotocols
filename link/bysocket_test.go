package link

import (
	"fmt"
	"log"
)

func ExampleNewBySocket() {
	peers := map[int]string{
		0: ":10000",
		1: ":10001",
	}
	l0, err := NewBySocket(0, "tcp4", peers)
	if err != nil {
		log.Fatalf("Error creation l0: %v", err)
	}
	l1, err := NewBySocket(1, "tcp4", peers)
	if err != nil {
		log.Fatalf("Error creation l1: %v", err)
	}
	// time.Sleep(time.Second)
	err = l0.Send(1, []byte("teste"))
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
	c := l1.GetDeliver()
	msg := <-c
	fmt.Println(msg.Src)
	fmt.Println(string(msg.Payload))
	// Output:
	// 0
	// teste
}
