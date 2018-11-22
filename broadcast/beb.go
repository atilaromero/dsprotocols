/*
DSPrototocols is a project that implements some of the protocols of
the book written by Christian Cachin, Rachid Guerraoui and Luís Rodrigues
"Introduction to Reliable and Secure Distributed Programming",
second edition, (https://www.distributedprogramming.net/), implemented
during the 2018 class "Tópicos especiais em processamento paralelo e
distribuído II" at Pontifícia Universidade Católica - RS, Brazil, under
supervision of professor Fernando Luis Dotti.

Copyright (C) 2018
	Atila Leites Romero (atilaromero@gmail.com),
	Carlos Renan Schick Louzada (crenan.louzada@gmail.com),
	Eliã Rafael de Lima Batista (o.elia.batista@gmail.com),
	Tarcisio Ceolin Junior (tarcisio@ceolin.org)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package broadcast

import (
	"log"

	"github.com/tarcisiocjr/dsprotocols/link"
)

// BebBroadcastMsg is a message to be broadcasted to all processes.
type BebBroadcastMsg struct {
	Payload []byte
}

// BebDelivertMsg contains the received brodcast message and the ID of the source process.
type BebDelivertMsg struct {
	Src     int
	Payload []byte
}

// Beb (best effort broadcast) is a struct that contains:
// Numproc: number of known processes.
// Pl: lower level perfect link
// Req: receives beb requests
// Ind: deliver beb messages
type Beb struct {
	NumProc int
	Pl      link.Link
	Req     chan BebBroadcastMsg
	Ind     chan BebDelivertMsg
}

// NewBeb returns a Beb struct, which implements Best Effort Broadcast.
// The sender also receives a copy of broadcasted messages.
//
// There are 4 channels here: the 2 channels from the perfect link
// are used for inter process communication, while req and ind are used to
// start and finish the broadcast.
//
// New broadcasts are initiated sending a message to the req channel.
//   Ex: beb.Req <- msg
//
// When an ongoing broadcast is received from another process,
// a deliver is triggered through the ind channel.
//   Ex: msg <- beb.Ind
//
// When using Beb, remember to create a go routine reading from the end channel beb.Ind and
// treating incomming messages.
func NewBeb(pl link.Link, numproc int) Beb {
	req := make(chan BebBroadcastMsg)
	ind := make(chan BebDelivertMsg)
	beb := Beb{numproc, pl, req, ind}

	go func() {
		plInd := pl.GetDeliver()
		for {
			select {
			// on new broadcast request
			case msg, ok := <-beb.Req:
				if !ok {
					return
				}
				// send the message to all known processes through each one's perfect link
				for q := 0; q < numproc; q++ {
					err := pl.Send(q, msg.Payload)
					if err != nil {
						log.Fatal(err)
					}
				}
			// when receiving a broadcast from another process
			case msg, ok := <-plInd:
				if !ok {
					return
				}
				// deliver the message one layer up
				go func() {
					beb.Ind <- BebDelivertMsg{
						Src:     msg.Src,
						Payload: msg.Payload,
					}
				}()
			}
		}
	}()

	return beb
}
