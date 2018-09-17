package consensus

import(
	"github.com/tarcisiocjr/dsprotocols/broadcast"
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

// EcBroadcastMsg is a message to be broadcasted to all processes.
type EcBroadcastMsg struct {
	Payload []byte
}

// EcDelivertMsg contains the received brodcast message and the ID of the current process.
type EcDelivertMsg struct {
	ID      int
	Payload []byte
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
	NumProc int
	Beb     broadcast.Beb
	Req     chan EcBroadcastMsg
	Ind     chan EcDelivertMsg
	Trusted	int
	Lastts	int
	Ts		int
}

func NewEc(beb broadcast.Beb, numproc int, trusted int, lastts int, ts int) Ec {
	req := make(chan EcBroadcastMsg)
	ind := make(chan EcDelivertMsg)
	ec := Ec{numproc, beb, req, ind, trusted, lastts, ts}

	bebs := []broadcast.Beb{}

	go func() {
		// on new ec request
		for msg, ok := <-ec.Req; ok; msg, ok = <-ec.Req {
			if trusted == 0 { // buscamos o id do processo em pl?
				//ts := ts + numproc
				for q := 0; q < numproc; q++{
					bebs[0].Req <- broadcast.BebBroadcastMsg{
						Payload: msg.Payload,
					}
				}
			}
		}
	}()

	return ec
}
