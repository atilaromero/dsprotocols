package consensus

import (
	"fmt"

	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/link"
)

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
	accepted int
	imLeader bool
}

func NewEp(pl link.Link, beb broadcast.Beb, totproc int) *Ep {

	accepted := 0
	tempval := 0
	state := State{0, 0}

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
		true,
	}

	// upon event ⟨ ep, Propose | v ⟩ do
	go func() {
		for msg, ok := <-ep.Req; ok; msg, ok = <-ep.Req {
			if !ep.imLeader {
				continue
			}
			ep.Tempval = msg.Val
			ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte("READ")}
		}
	}()

	// upon event ⟨ beb, Deliver | l, [READ] ⟩ do
	go func() {
		for msg, ok := <-beb.Ind; ok; msg, ok = <-beb.Ind {
			pl.Send(msg.Src, []byte(fmt.Sprintf("[STATE,%d,%d]\n", ep.State.ValTS, ep.State.Val)))
			ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte("READ")}
		}
	}()

	// upon event ⟨ pl, Deliver | q, [STATE, ts, v] ⟩ do
	go func() {
		ind := pl.GetDeliver()
		for msg, ok := <-ind; ok; msg, ok = <-ind {
			if !ep.imLeader {
				continue
			}
			var ts int
			var v int
			fmt.Sscanf(string(msg.Payload), "[STATE,%d,%d]\n", &ts, &v)
			ep.States[msg.Src] = State{ValTS: ts, Val: v}
		}
	}()

	return &ep
}
