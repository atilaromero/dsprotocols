package main

import "fmt"

type Beb struct {
	ID        int
	UpLayer   P2PLink
	SameLayer P2PLink
	Contacts  []Beb
}

func NewBeb(id int) Beb {
	result := Beb{
		ID: id,
	}
	result.UpLayer = NewP2PLink(NewCommunicationHandler(func(msg []byte) {
		for _, c := range result.Contacts {
			c.SameLayer.Req <- msg
		}
	}, func(msg []byte) {
		fmt.Println("all messages delivered")
	}))
	result.SameLayer = NewP2PLink(NewCommunicationHandler(func(msg []byte) {

	}, func(msg []byte) {

	}))
	return result
}
