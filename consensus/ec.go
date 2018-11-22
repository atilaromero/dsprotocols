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
package consensus

import (
	"fmt"
	"log"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/leadership"
	"github.com/tarcisiocjr/dsprotocols/link"
)

/*
	Properties of Epoch-change
		EC1: Monotonicity: If a correct process starts an epoch (ts,l) and later starts
			an epoch (ts',l'), then ts' > ts.
		EC2: Consistency: If a correct process starts an epoch (ts,l) and another correct
			process starts an epoch (ts',l') with ts = ts', then l = l'.
		EC3: Eventual leadership: There is a time after which every correct process has
			started some epoch and starts no further epoch, such that the last epoch
			started at every correct process is epoch (ts,l) and process l is correct.
*/

/*
	Every process p maintains two timestamps: a timestamp lastts of the last epoch that is
	started (i.d., for which it triggered a <StartEpoch> event), and the timestamp ts of
	the last epoch that is attempted to start with itself as leader (i.e., for which it
	broadcast a NEWEPOCH message).
	Initially, the process sets ts to its rank. Whenever the leader detector subsequently
	makes p trust itself, p adds N to ts and sends a NEWEPOCH message with ts.
*/

// EcDelivertMsg contains the received brodcast message and the ID of the current process.
type EcDelivertMsg struct {
	Ts     int
	Leader int
}

// Ec (Epoch-change) is a struct that contains:
// Numproc: number of known processes.
// Pl: lower level perfect link
// Req: receives beb requests
// Ind: deliver beb messages
// Trusted:
// Lastts: last epoch started by the process
// Ts: last timestamp attempted to start by the process
type Ec struct {
	Pl             link.Link
	Beb            broadcast.Beb
	LeaderDetector <-chan leadership.TrustMsg
	Ind            chan EcDelivertMsg
	TotProc        int
	Trusted        int
	Lastts         int
	Ts             int
}

func NewEc(pl link.Link, beb broadcast.Beb, omega <-chan leadership.TrustMsg, totproc int) *Ec {

	//upon event < ec, Init > do
	trusted := 0
	lastts := 0
	ts := pl.ID()

	ec := Ec{pl, beb, omega, make(chan EcDelivertMsg), totproc, trusted, lastts, ts}

	// upon event < Ω , Trust | p > do
	onOmega := func(p leadership.TrustMsg) {
		// from book's errata
		if p.ID != ec.Trusted {
			err := ec.Pl.Send(ec.Trusted, []byte("NACK"))
			if err != nil {
				log.Fatal(err)
			}
		}
		ec.Trusted = p.ID
		if ec.Trusted == pl.ID() {
			ec.Ts += ec.TotProc
			beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("%d", ec.Ts))}
		}
	}

	// upon event < beb, Deliver | l , [ NEWEPOCH , newts ] > do
	onBeb := func(msg broadcast.BebDelivertMsg) {
		newts := 0
		_, err := fmt.Sscanf(string(msg.Payload), "%d", &newts)
		if err != nil {
			log.Fatal(err)
		}
		if msg.Src == ec.Trusted && newts > ec.Lastts {
			ec.Lastts = newts
			ec.Ind <- EcDelivertMsg{Ts: newts, Leader: msg.Src}
		} else {
			err := ec.Pl.Send(msg.Src, []byte("NACK"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// upon event < pl, Deliver | p , [ NACK ] > do
	onPl := func(msg link.Message) {
		if ec.Trusted == pl.ID() {
			ec.Ts += ec.TotProc
			beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("%d", ec.Ts))}
		}
	}

	go func() {
		plInd := pl.GetDeliver()
		defer close(ec.Ind)
		for {
			select {
			case p, ok := <-ec.LeaderDetector:
				if !ok {
					return
				}
				onOmega(p)
			case msg, ok := <-ec.Beb.Ind:
				if !ok {
					return
				}
				onBeb(msg)
			case msg, ok := <-plInd:
				if !ok {
					return
				}
				onPl(msg)
			}
		}
	}()
	return &ec
}
