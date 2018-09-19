package consensus

import (
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
	trusted := 0
	lastts := 0
	ts := pl.ID()
	ec := Ec{pl, beb, omega, make(chan EcDelivertMsg), totproc, trusted, lastts, ts}

	// upon event < Ω , Trust | p > do
	go func() {
		for p, ok := <-ec.LeaderDetector; ok; p, ok = <-ec.LeaderDetector {
			if p.ID != ec.Trusted {
				err := ec.Pl.Send(ec.Trusted, []byte("NACK"))
				if err != nil {
					log.Fatal(err)
				}
			}
			ec.Trusted = p.ID
			if ec.Trusted == pl.ID() {
				ec.Ts += ec.TotProc
				beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte{byte(ec.Ts)}}
			}
		}
	}()

	// upon event < beb, Deliver | l , [ NEWEPOCH , newts ] > do
	go func() {
		for msg, ok := <-ec.Beb.Ind; ok; msg, ok = <-ec.Beb.Ind {
			newts := int(msg.Payload[0])
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
	}()

	go func() {
		plInd := pl.GetDeliver()
		for _, ok := <-plInd; ok; _, ok = <-plInd {
			if ec.Trusted == pl.ID() {
				ec.Ts += ec.TotProc
				beb.Req <- broadcast.BebBroadcastMsg{Payload: []byte{byte(ec.Ts)}}
			}
		}
	}()

	return &ec
}
