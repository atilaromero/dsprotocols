
# link
--
    import "github.com/tarcisiocjr/dsprotocols/link"


## Usage

#### type Link

```go
type Link interface {
	Send(id int, payload []byte) error
	GetDeliver() <-chan Message
	ID() int
}
```

Link is a struct that contains the ID of a process and 2 channels: Link receives
messages from a upper layer through Req and deliver messages through Ind.

#### func  NewByChan

```go
func NewByChan(id int, peers map[int]chan<- Message) Link
```
NewByChan returns a new Pl struct. This implementations uses only channels,
instead of sockets, designed to be used in a non distributed environment (single
machine).

For a network implementation with similar behavior, see XXX...XXX

ID refers to the process than owns this Pl. peers is a, possibly empty, pool of
all Pl. NewPl will modify this pool, incluing the new Pl in it.

The requisition channel is used to receive messages from the upper layer.

    Ex: pl.Req <- msg

When a new message is received through Req, that message's destination is seeked
in the pool of perfect links (peers), and sent through the indication channel.
Is up to the upper layer to collect messages from the indication channel.

    Ex: msg <- pl.Ind

The indication channel will block if those messages aren't treated somehow. To
prevent that, create a go routine that continually reads from pl.Ind

#### func  NewBySocket

```go
func NewBySocket(id int, proto string, peers map[int]string) (Link, error)
```
NewBySocket returns a new `Socket` or an `error`. NOTE: All connections to
`Peer`s are established during this function.

#### type Message

```go
type Message struct {
	Src     int
	Payload []byte
}
```

Message is use for indication. These messages are intended to be read by the
upper layer. Src is the sender's process ID in the link pool.

# broadcast
--
    import "github.com/tarcisiocjr/dsprotocols/broadcast"


## Usage

#### type Beb

```go
type Beb struct {
	NumProc int
	Pl      link.Link
	Req     chan BebBroadcastMsg
	Ind     chan BebDelivertMsg
}
```

Beb (best effort broadcast) is a struct that contains: Numproc: number of known
processes. Pl: lower level perfect link Req: receives beb requests Ind: deliver
beb messages

#### func  NewBeb

```go
func NewBeb(pl link.Link, numproc int) Beb
```
NewBeb returns a Beb struct, which implements Best Effort Broadcast. The sender
also receives a copy of broadcasted messages.

There are 4 channels here: the 2 channels from the perfect link are used for
inter process communication, while req and ind are used to start and finish the
broadcast.

New broadcasts are initiated sending a message to the req channel.

    Ex: beb.Req <- msg

When an ongoing broadcast is received from another process, a deliver is
triggered through the ind channel.

    Ex: msg <- beb.Ind

When using Beb, remember to create a go routine reading from the end channel
beb.Ind and treating incomming messages.

#### type BebBroadcastMsg

```go
type BebBroadcastMsg struct {
	Payload []byte
}
```

BebBroadcastMsg is a message to be broadcasted to all processes.

#### type BebDelivertMsg

```go
type BebDelivertMsg struct {
	Src     int
	Payload []byte
}
```

BebDelivertMsg contains the received brodcast message and the ID of the source
process.

# consensus
--
    import "github.com/tarcisiocjr/dsprotocols/consensus"


## Usage

#### type Ec

```go
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
```

Ec (Epoch-change) is a struct that contains: Numproc: number of known processes.
Pl: lower level perfect link Req: receives beb requests Ind: deliver beb
messages Trusted: Lastts: last epoch started by the process Ts: last timestamp
attempted to start by the process

#### func  NewEc

```go
func NewEc(pl link.Link, beb broadcast.Beb, omega <-chan leadership.TrustMsg, totproc int) *Ec
```

#### type EcDelivertMsg

```go
type EcDelivertMsg struct {
	Ts     int
	Leader int
}
```

EcDelivertMsg contains the received brodcast message and the ID of the current
process.

#### type Ep

```go
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
```

Ep (Epoch-consensus) is a struct that contains: Pl: lower level perfect link
Beb: receives and sends beb requests Ind: deliver EpDecideMsg messages Req:
receives EpProposeMsg messages TotProc: number of known processes

#### func  NewEp

```go
func NewEp(pl link.Link, beb broadcast.Beb, totproc int) *Ep
```

#### func (*Ep) Init

```go
func (ep *Ep) Init(leader int, pState State)
```

#### type EpDecideMsg

```go
type EpDecideMsg struct {
	Abort bool
	State State
}
```

EpDecideMsg contains the decided value

#### type EpProposeMsg

```go
type EpProposeMsg struct {
	Abort bool
	Val   int
}
```

EpProposeMsg contains the value to be agreed

#### type State

```go
type State struct {
	ValTS int
	Val   int
}
```


#### type Uc

```go
type Uc struct {
	Req        chan UcProposeMsg
	Ind        chan UcDeliverMsg
	EcInstance *Ec
	EpInstance *Ep
	Val        int
	Proposed   bool
	Decided    bool
	Ets        int
	L          int
	NewTS      int
	NewL       int
}
```

Uc (Uniform Consensus,) is a struct that contains: EcInstance: epoch change
instance map of ep: multiple instances of epoch consensus

#### func  NewUC

```go
func NewUC(ec *Ec, ep *Ep) *Uc
```

#### func (*Uc) Init

```go
func (uc *Uc) Init()
```

#### type UcDeliverMsg

```go
type UcDeliverMsg struct {
	Val int
}
```


#### type UcProposeMsg

```go
type UcProposeMsg struct {
	Val int
}
```

# leadership
--
    import "github.com/tarcisiocjr/dsprotocols/leadership"


## Usage

#### type LeaderDetector

```go
type LeaderDetector struct {
	Suspected []int
	Leader    int
	Ind       chan TrustMsg
}
```

LeaderDetector (Eventual Leader Detector) is a struct that contains: Ind:
deliver TrustMsg indicating the new leader.

#### func  NewLeaderDetector

```go
func NewLeaderDetector(numProc int) LeaderDetector
```
NewLeaderDetector returns a NewLeaderDetector struct, which implements (fake)
Eventual Leader Detector.

There is 1 channel here: ind is used to indicate a leader to the upper layer.
Receives the number of processes

#### type TrustMsg

```go
type TrustMsg struct {
	ID int
}
```

TrustMsg contains the ID of the current leader process.
