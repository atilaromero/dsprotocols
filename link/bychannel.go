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
package link

type linkChannel struct {
	id         int
	indication chan Message
	peers      map[int]chan<- Message
}

func (l linkChannel) Send(id int, payload []byte) error {
	// find destination in pool and fill Src with current process ID
	go func() {
		// fmt.Printf("SEND src:%d dst:%d '%s'\n", l.id, id, string(payload))
		l.peers[id] <- Message{
			Src:     l.id,
			Payload: payload,
		}
	}()
	return nil
}

func (l linkChannel) GetDeliver() <-chan Message {
	return l.indication
}

func (l linkChannel) ID() int {
	return l.id
}

// NewByChan returns a new Pl struct. This implementations uses
// only channels, instead of sockets, designed to be used in a non distributed environment (single machine).
//
// For a network implementation with similar behavior, see XXX...XXX
//
// ID refers to the process than owns this Pl.
// peers is a, possibly empty, pool of all Pl. NewPl will modify this pool, incluing the new Pl in it.
//
// The requisition channel is used to receive messages from the upper layer.
//   Ex: pl.Req <- msg
// When a new message is received through Req,
// that message's destination is seeked in the pool of perfect links (peers), and sent through the
// indication channel.
// Is up to the upper layer to collect messages from the indication channel.
//   Ex: msg <- pl.Ind
//
// The indication channel will block if those messages aren't treated somehow.
// To prevent that, create a go routine that continually reads from pl.Ind
func NewByChan(id int, peers map[int]chan<- Message) Link {
	ind := make(chan Message)
	l := linkChannel{id, ind, peers}

	// adds the new created Pl instance to the pool of peers
	peers[id] = ind

	return l
}
