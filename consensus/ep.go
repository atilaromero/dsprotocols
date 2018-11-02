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
	Leader   int
	Ets      int
	Aborted  chan bool
}

func highest(states map[int]State) State {
	maxts := State{Val: -1, ValTS: -1}
	for _, state := range states {
		if state.ValTS > maxts.ValTS {
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
	}

	return &ep
}

func (ep *Ep) Init(ets int, leader int) {

	ep.State = State{Val: -1, ValTS: -1}
	ep.Tempval = -1
	ep.States = make(map[int]State)
	ep.Accepted = 0
	ep.Leader = leader
	ep.Ets = ets
	ep.Aborted = make(chan bool)

	// upon event ⟨ ep, Propose | v ⟩ do
	// OR
	// upon event ⟨ ep, Abort ⟩ do
	go func() {
		for msg, ok := <-ep.Req; ok; msg, ok = <-ep.Req {

			// when receiving abort message
			if msg.Abort {
				comment("Aborting...")
				go func() {
					ep.Ind <- EpDecideMsg{Abort: true, AbortedState: ep.State}
				}()
				close(ep.Aborted)
				return
			}

			// only leader l
			if ep.Leader != ep.Pl.ID() {
				continue
			}

			// when leader received Propose message
			ep.Tempval = msg.Val
			comment("Broadcasting READ")
			go func() {
				ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte("READ 0")} //0 is a dummy value to make all 3 messages (READ, WRITE, and DECIDED) have the same format
			}()
		}
	}()

	// upon event ⟨ beb, Deliver | l, [READ] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [WRITE, v] ⟩ do
	// OR
	// upon event ⟨ beb, Deliver | l, [DECIDED, v] ⟩ do
	go func() {
		for {
			select {
			case <-ep.Aborted:
				return
			case msg, ok := <-ep.Beb.Ind:
				if !ok {
					panic(nil)
				}

				v := -1
				msgType := ""
				_, err := fmt.Sscanf(string(msg.Payload), "%s %d", &msgType, &v)
				if err != nil {
					log.Panic(err, ";Payload: ", string(msg.Payload), msg.Payload)
				}

				switch msgType {
				case "READ":
					comment("Sending STATE")
					ep.Pl.Send(msg.Src, []byte(fmt.Sprintf("STATE %d %d", ep.State.ValTS, ep.State.Val)))
					continue
				case "WRITE":
					ep.State = State{ValTS: ep.Ets, Val: v}
					comment(fmt.Sprintf("Received [WRITE,%d] ", v))
					ep.Pl.Send(msg.Src, []byte("ACCEPT"))
					continue
				case "DECIDED":
					comment(fmt.Sprintf("Received [DECIDED,%d]", v))
					go func() {
						ep.Ind <- EpDecideMsg{Abort: false, Val: v, AbortedState: State{}}
					}()
				default:
					log.Panic("unknow msgType", msgType)
				}
			}
		}
	}()

	// upon event ⟨ pl, Deliver | q, [STATE, ts, v] ⟩ do
	// OR
	// upon event ⟨ pl, Deliver | q , [ACCEPT] ⟩
	go func() {
		ind := ep.Pl.GetDeliver()
		for {
			select {
			case <-ep.Aborted:
				return
			case msg, ok := <-ind:
				if !ok {
					return
				}
				// only leader l
				if ep.Leader != ep.Pl.ID() {
					continue
				}

				var ts int
				var v int

				// if its an ACCEPT message
				if string(msg.Payload) == "ACCEPT" {
					ep.Accepted = ep.Accepted + 1
					comment(fmt.Sprintf("Received ACCEPTS: %d", ep.Accepted))
					//upon accepted > N/2 do
					if ep.Accepted > ep.TotProc/2 {
						ep.Accepted = 0
						comment("Broadcasting DECIDED")
						go func() {
							ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("DECIDED %d", ep.Tempval))}
						}()
					}
					continue
				}

				// otherwise is a STATE message
				_, err := fmt.Sscanf(string(msg.Payload), "STATE %d %d", &ts, &v)
				if err != nil {
					panic(err)
				}
				comment(fmt.Sprintf("Received [STATE,%d,%d] from process %d", ts, v, msg.Src))
				ep.States[msg.Src] = State{ValTS: ts, Val: v}
				// upon #( states ) > N/2 do
				if len(ep.States) > ep.TotProc/2 {
					state := highest(ep.States)
					if state.Val != -1 {
						ep.Tempval = state.Val
					}
					ep.States = make(map[int]State)
					comment("Broadcasting WRITE", ep.Tempval)
					go func() {
						ep.Beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte(fmt.Sprintf("WRITE %d", ep.Tempval))}
					}()
				}
			}
		}
	}()
}
