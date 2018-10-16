package consensus

import (
	"fmt"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

const printLog bool = true

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

/*
	Description...
*/

// EpProposeMsg contains the value to be agreed
type EpProposeMsg struct {
	Abort bool
	Val   int
}

// EpDecideMsg contains the decided value
type EpDecideMsg struct {
	Abort        bool
	Val          int
	AbortedState State
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
	IsLeader bool
}

func NewEp(pl link.Link, beb broadcast.Beb, totproc int, pState State, pIsLeader bool) *Ep {

	accepted := 0
	tempval := 0
	state := pState

	ep := Ep{
		pl,
		beb,
		make(chan EpDecideMsg),
		make(chan EpProposeMsg),
		totproc,
		state,
		tempval,
		make(map[int]State),
		accepted,
		pIsLeader,
	}

	// upon event ⟨ ep, Propose | v ⟩ do
	go func() {
		for msg, ok := <-ep.Req; ok; msg, ok = <-ep.Req {
			if !ep.IsLeader {
				continue
			}
			ep.Tempval = msg.Val
			if printLog {
				fmt.Println("Broadcasting READ")
			}
			ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte("READ")}
		}
	}()

	// upon event ⟨ beb, Deliver | l, [READ] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [WRITE, v] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [DECIDED, v] ⟩ do
	go func() {
		for msg, ok := <-beb.Ind; ok; msg, ok = <-beb.Ind {

			var v int
			var ets int

			ets = ep.State.ValTS // ????
			v = -1

			// if its a READ message
			if string(msg.Payload) == "READ" {
				if printLog {
					fmt.Println("Sending STATE")
				}
				pl.Send(msg.Src, []byte(fmt.Sprintf("[STATE,%d,%d]\n", ep.State.ValTS, ep.State.Val)))
				continue
			}

			fmt.Sscanf(string(msg.Payload), "[WRITE,%d]\n", &v)

			// else if its a WRITE message
			if v != -1 {
				ep.State = State{ValTS: ets, Val: v}
				if printLog {
					fmt.Println(fmt.Sprintf("Received [WRITE,%d] ", v))
				}
				pl.Send(msg.Src, []byte("ACCEPT"))
				continue
			}

			fmt.Sscanf(string(msg.Payload), "[DECIDED,%d]\n", &v)

			// otherwise its a DECIDE message
			if printLog {
				fmt.Println(fmt.Sprintf("Received [DECIDED,%d]", v))
			}
			ep.Ind <- EpDecideMsg{Abort: false, Val: v, AbortedState: State{}}
		}
	}()

	// upon event ⟨ pl, Deliver | q, [STATE, ts, v] ⟩ do
	// OR
	// upon event ⟨ pl, Deliver | q , [ACCEPT] ⟩
	go func() {
		ind := pl.GetDeliver()
		for msg, ok := <-ind; ok; msg, ok = <-ind {

			if !ep.IsLeader {
				continue
			}

			var ts int
			var v int

			// if its an ACCEPT message
			if string(msg.Payload) == "ACCEPT" {
				ep.Accepted = ep.Accepted + 1
				if printLog {
					fmt.Println(fmt.Sprintf("Received ACCEPTS: %d", ep.Accepted))
				}
				continue
			}

			// otherwise is a STATE message
			fmt.Sscanf(string(msg.Payload), "[STATE,%d,%d]\n", &ts, &v)

			if printLog {
				fmt.Println(fmt.Sprintf("Received [STATE,%d,%d] from process %d", ts, v, msg.Src))
			}
			ep.States[msg.Src] = State{ValTS: ts, Val: v}
		}
	}()

	go func() {
		for {
			// upon #( states ) > N/2 do
			if len(ep.States) > ep.TotProc/2 {

				// TODO:
				/*
					(ts, v) := highest ( states ) ;
					if v = ⊥ then
						tmpval := v ;
					states := [⊥] N ;
					trigger  beb, Broadcast | [ W RITE , tmpval ]  ;
				*/

				ep.Tempval = 1 // pegar o maior valor do states

				ep.States = make(map[int]State)
				if printLog {
					fmt.Println("Broadcasting WRITE")
				}
				ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("[%s,%d]\n", "WRITE", ep.Tempval))}
				break
			}
		}

		//upon accepted > N/2 do
		for {
			if ep.Accepted > ep.TotProc/2 {
				ep.Accepted = 0
				if printLog {
					fmt.Println("Broadcasting DECIDED")
				}
				ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("[DECIDED,%d]\n", ep.Tempval))}
				break
			}
		}
	}()

	return &ep
}
