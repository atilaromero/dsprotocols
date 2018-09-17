package link

import "fmt"

func ExampleNewByChan() {
	peers := make(map[int]chan<- Message)
	l0 := NewByChan(0, peers)
	l1 := NewByChan(1, peers)
	l0.Send(1, []byte("teste"))
	c := l1.GetDeliver()
	msg := <-c
	fmt.Println(msg.Src)
	fmt.Println(string(msg.Payload))
	// Output:
	// 0
	// teste
}
