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

// Message is use for indication.
// These messages are intended to be read by the upper layer.
// Src is the sender's process ID in the link pool.
type Message struct {
	Src     int
	Payload []byte
}

// Link is a struct that contains the ID of a process and 2 channels:
// Link receives messages from a upper layer through Req
// and deliver messages through Ind.
type Link interface {
	Send(id int, payload []byte) error
	GetDeliver() <-chan Message
	ID() int
}
