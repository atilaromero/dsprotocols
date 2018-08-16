package try2

import (
	"sync"
)

type process struct {
	id int
	c  chan []byte
}

func newProcess(id int) process {
	return process{
		id,
		make(chan []byte),
	}
}

type beb2 struct {
	channels []chan []byte
}

func (b *beb2) send(msg []byte) {
	wg := new(sync.WaitGroup)
	for _, c := range b.channels {
		wg.Add(1)
		go func(c chan []byte) {
			c <- msg
			defer wg.Done()
		}(c)
	}
	wg.Wait()
}
