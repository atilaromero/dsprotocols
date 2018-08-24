package layer

import (
	"fmt"
	"time"
)

func ExampleNew() {
	pool := map[int]Layer{}
	l0 := New(0, pool)
	l0.UpponEvent("SEND", func(msg Msg) {
		fmt.Println(string(msg.Payload))
	})
	l0.Trigger(0, "SEND", []byte("teste"))
	time.Sleep(1 * time.Second)
	// Output:
	// teste
}
