package main

import "testing"

func TestA(t *testing.T) {
	msg := []byte("hello world")
	p1 := newProcess(1)
	p2 := newProcess(2)
	if p1 == p2 {
		t.Fail()
	}
	b := beb2{
		channels: []chan []byte{p1.c, p2.c},
	}
	go func() {
		b.send(msg)
	}()
	msg1 := <-p1.c
	msg2 := <-p2.c
	if string(msg1) != string(msg) {
		t.Fail()
	}
	if string(msg2) != string(msg) {
		t.Fail()
	}
}
