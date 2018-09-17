# link
--
    import "github.com/tarcisiocjr/dsprotocols/link"


## Usage

#### type Pl

```go
type Pl struct {
	ID  int
	Req chan PlSendMsg
	Ind chan PlDeliverMsg
}
```

Pl is a struct that contains the ID of a process and 2 channels: Pl receives
messages from a upper layer through Req and deliver messages through Ind.

#### func  NewPl

```go
func NewPl(ID int, pls map[int]Pl) Pl
```
NewPl returns a new Pl struct. This implementations uses only channels, instead
of sockets, designed to be used in a non distributed environment (single
machine).

For a network implementation with similar behavior, see XXX...XXX

ID refers to the process than owns this Pl. pls is a, possibly empty, pool of
all Pl. NewPl will modify this pool, incluing the new Pl in it.

The requisition channel is used to receive messages from the upper layer.

    Ex: pl.Req <- msg

When a new message is received through Req, that message's destination is seeked
in the pool of perfect links (pls), and sent through the indication channel. Is
up to the upper layer to collect messages from the indication channel.

    Ex: msg <- pl.Ind

The indication channel will block if those messages aren't treated somehow. To
prevent that, create a go routine that continually reads from pl.Ind

#### type PlDeliverMsg

```go
type PlDeliverMsg struct {
	Src     int
	Payload []byte
}
```

PlDeliverMsg is use for indication. These messages are intended to be read by
the upper layer. Src is the sender's process ID in the link pool.

#### type PlSendMsg

```go
type PlSendMsg struct {
	Dst     int
	Payload []byte
}
```

PlSendMsg is used for requests. These messages are intended to be created by the
upper layer and sent to Pl. Dst is the destination process ID in the link pool.


# broadcast
--
    import "github.com/tarcisiocjr/dsprotocols/broadcast"


## Usage

#### type Beb

```go
type Beb struct {
	NumProc int
	Pl      link.Pl
	Req     chan BebBroadcastMsg
	Ind     chan BebDelivertMsg
}
```

Beb (best effort broadcast) is a struct that contains: Numproc: number of known
processes. Pl: lower level perfect link Req: receives beb requests Ind: deliver
beb messages

#### func  NewBeb

```go
func NewBeb(pl link.Pl, numproc int) Beb
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
	ID      int
	Payload []byte
}
```

BebDelivertMsg contains the received brodcast message and the ID of the current
process.
