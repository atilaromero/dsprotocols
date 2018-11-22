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
	"github.com/tarcisiocjr/dsprotocols/link"
)

func comment(a ...interface{}) {
	fmt.Println(a...)
}

/*
	Properties of Epoch-consensus
		EP1: Validity: If a correct process ep-decides v, then v was ep-proposed by the
		leader l′ of some epoch consensus with timestamp ts′ ≤ ts and leader l′.

		EP2: Uniform agreement: No two processes ep-decide differently.

		EP3: Integrity: Every correct process ep-decides at most once.

		EP4: Lock-in: If a correct process has ep-decided v in an epoch consensus with timestamp
		ts′ < ts, then no correct process ep-decides a value different from v.

		EP5: Termination: If the leader l is correct, has ep-proposed a value, and no
		correct process aborts this epoch consensus, then every correct process eventually
		ep-decides some value.

		EP6: Abort behavior: When a correct process aborts an epoch consensus, it eventually
		will have completed the abort; moreover, a correct process completes an abort only if
		the epoch consensus has been aborted by some correct process.
*/

// EpProposeMsg contains the value to be agreed
type EpProposeMsg struct {
	Abort bool
	Val   int
}

// EpDecideMsg contains the decided value
type EpDecideMsg struct {
	Abort bool
	State State
}

type State struct {
	ValTS int
	Val   int
}

// Ep (Epoch-consensus) is a struct that contains:
// Pl: lower level perfect link
// Beb: receives and sends beb requests
// Ind: deliver EpDecideMsg messages
// Req: receives EpProposeMsg messages
// TotProc: number of known processes
type Ep struct {
	Pl       link.Link
	Beb      broadcast.Beb
	Ind      chan EpDecideMsg
	Req      chan EpProposeMsg
	TotProc  int
	State    State
	Tempval  int
	States   map[int]State
	Accepted int
	Leader   int
}

func highest(states map[int]State) State {
	maxts := State{Val: -1, ValTS: -1}
	for _, state := range states {
		if state.Val != -1 && state.ValTS > maxts.ValTS {
			maxts = state
		}
	}
	return maxts
}

func NewEp(pl link.Link, beb broadcast.Beb, totproc int) *Ep {

	ep := Ep{
		Pl:      pl,
		Beb:     beb,
		Ind:     make(chan EpDecideMsg),
		Req:     make(chan EpProposeMsg),
		TotProc: totproc,
		State:   State{ValTS: -1, Val: -1},
	}

	return &ep
}

func (ep *Ep) Init(leader int, pState State) {

	// upon event < ep, Init | state > do
	// 		( valts , val ) := state;
	// 		tmpval := ⊥ ;
	// 		states := [⊥] N ;
	// 		accepted := 0 ;

	ep.State.ValTS = pState.ValTS
	ep.Tempval = -1
	ep.States = make(map[int]State)
	ep.Accepted = 0
	ep.Leader = leader

	aborted := make(chan bool)

	// upon event ⟨ ep, Propose | v ⟩ do
	// OR
	// upon event ⟨ ep, Abort ⟩ do
	onRequest := func(msg EpProposeMsg, aborted chan bool) (ok bool) {
		// upon event < ep, Abort > do
		// 		trigger < ep, Aborted | ( valts , val ) > ;
		// 		halt;
		if msg.Abort {
			comment("Aborting...")
			go func(msg EpDecideMsg) {
				ep.Ind <- msg
			}(EpDecideMsg{Abort: true, State: ep.State})
			close(aborted)
			return false
		}

		// only leader l
		if ep.Leader != ep.Pl.ID() {
			return true
		}

		// upon event ⟨ ep, Propose | v ⟩ do
		// 		tmpval := v ;
		//		trigger < beb, Broadcast | [ READ ] > ;
		ep.Tempval = msg.Val
		comment("Broadcasting READ")
		go func(msg broadcast.BebBroadcastMsg) {
			ep.Beb.Req <- msg
			//0 is a dummy value to make all 3 messages (READ, WRITE, and DECIDED) have the same format
		}(broadcast.BebBroadcastMsg{Payload: []byte("READ 0 0")})
		return true
	}

	// upon event ⟨ beb, Deliver | l, [READ] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [WRITE, v] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [DECIDED, v] ⟩ do
	onBeb := func(msg broadcast.BebDelivertMsg) {
		v := -1
		vts := -1
		msgType := ""
		_, err := fmt.Sscanf(string(msg.Payload), "%s %d %d", &msgType, &vts, &v)
		if err != nil {
			log.Panic(err, ";Payload: ", string(msg.Payload), msg.Payload)
		}

		switch msgType {
		// 	upon event < beb, Deliver | l , [ READ ] > do
		// 		trigger < pl, Send | l , [ STATE , valts , val ] > ;
		case "READ":
			comment("Sending STATE")
			ep.Pl.Send(msg.Src, []byte(fmt.Sprintf("STATE %d %d", ep.State.ValTS, ep.State.Val)))
			return
		// upon event < beb, Deliver | l , [ WRITE , v] > do
		// 		( valts , val ) := (ets, v) ;
		// 		trigger < pl, Send | l , [ ACCEPT ] > ;
		case "WRITE":
			// ignore old ets
			if vts < ep.State.ValTS {
				return
			}
			ep.State = State{ValTS: vts, Val: v}
			comment(fmt.Sprintf("Received [WRITE,%d,%d] ", v, vts))
			ep.Pl.Send(msg.Src, []byte("ACCEPT"))
			return
		// upon event < beb, Deliver | l , [ DECIDED , v] > do
		// 		trigger < ep, Decide | v > ;
		case "DECIDED":
			// ignore old ets
			if vts < ep.State.ValTS {
				return
			}
			comment(fmt.Sprintf("Received [DECIDED,%d]", v))
			ep.State = State{ValTS: vts, Val: v}
			go func(msg EpDecideMsg) {
				ep.Ind <- msg
			}(EpDecideMsg{Abort: false, State: ep.State})
		default:
			log.Panic("unknown msgType", msgType)
		}
	}

	// upon event ⟨ pl, Deliver | q, [STATE, ts, v] ⟩ do
	// OR
	// upon event ⟨ pl, Deliver | q , [ACCEPT] ⟩
	onPl := func(msg link.Message) {
		// only leader l
		if ep.Leader != ep.Pl.ID() {
			return
		}

		// upon event < pl, Deliver | q , [ ACCEPT ] > do
		// 		accepted := accepted + 1 ;
		if string(msg.Payload) == "ACCEPT" {
			ep.Accepted = ep.Accepted + 1
			comment(fmt.Sprintf("Received ACCEPTS: %d", ep.Accepted))
			// upon accepted > N/2 do
			// 		accepted := 0 ;
			// 		trigger < beb, Broadcast | [ DECIDED , tmpval ] >  ;
			if ep.Accepted > ep.TotProc/2 {
				ep.Accepted = 0
				comment("Broadcasting DECIDED")
				go func(msg broadcast.BebBroadcastMsg) {
					ep.Beb.Req <- msg
				}(broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("DECIDED %d %d", ep.State.ValTS, ep.Tempval))})
			}
			return
		}

		// otherwise is a STATE message
		// upon event < pl, Deliver | q , [ STATE , ts, v] > do
		// 		states [q] := (ts, v) ;
		var ts int
		var v int
		_, err := fmt.Sscanf(string(msg.Payload), "STATE %d %d", &ts, &v)
		if err != nil {
			panic(err)
		}
		comment(fmt.Sprintf("Received [STATE,%d,%d] from process %d", ts, v, msg.Src))
		ep.States[msg.Src] = State{ValTS: ts, Val: v}

		// upon #( states ) > N/2 do
		// 		(ts, v) := highest ( states ) ;
		// 		if v = ⊥ then
		// 		tmpval := v ;
		// 		states := [⊥] N ;
		// 		trigger < beb, Broadcast | [ WRITE , tmpval ] > ;
		if len(ep.States) > ep.TotProc/2 {
			state := highest(ep.States)
			if state.Val != -1 {
				ep.Tempval = state.Val
			}
			ep.States = make(map[int]State)
			comment("Broadcasting WRITE", ep.Tempval)
			go func(msg broadcast.BebBroadcastMsg) {
				ep.Beb.Req <- msg
			}(broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("WRITE %d %d", ep.State.ValTS, ep.Tempval))})
		}
	}

	go func(ep *Ep, aborted chan bool) {
		ind := ep.Pl.GetDeliver()
		for {
			select {
			case <-aborted:
				return
			case msg, ok := <-ep.Req:
				if !ok {
					return
				}
				ok = onRequest(msg, aborted)
				if !ok {
					return
				}
			case msg, ok := <-ep.Beb.Ind:
				if !ok {
					return
				}
				onBeb(msg)

			case msg, ok := <-ind:
				if !ok {
					return
				}
				onPl(msg)
			}
		}
	}(ep, aborted)
}
