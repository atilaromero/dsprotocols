package consensus

import (
	"github.com/tarcisiocjr/dsprotocols/broadcast"
	"github.com/tarcisiocjr/dsprotocols/leadership"
	"github.com/tarcisiocjr/dsprotocols/link"
)

/*
	Algorithm 5.7: Leader-Driven Consensus
	Implements:
		UniformConsensus, instance uc.
	Uses:
		EpochChange, instance ec;
		EpochConsensus (multiple instances).
*/

// Uc (Uniform Consensus,) is a struct that contains:
// EcInstance: epoch change instance
// map of ep: multiple instances of epoch consensus
type Uc struct {
	EcInstance *Ec
	Consensus  map[int]Ep
	val        int
	validVal   bool
	proposed   bool
	decided    bool
	ets        int
	l          int
	newts      int
	newl       int
}

func NewUC(pl link.Link, beb broadcast.Beb, omega <-chan leadership.TrustMsg, totproc int) *Uc {

	//upon event ⟨ uc, Init ⟩ do
	// ep := NewEp(pl2, beb2, totproc, 0, l)
	uc := &Uc{
		EcInstance: NewEc(pl, beb, omega, totproc),
		Consensus:  make(map[int]Ep),
		validVal:   false,
		val:        0,
		proposed:   false,
		decided:    false,
		ets:        0,
		newts:      0,
		newl:       -1,
	}
	uc.l = uc.EcInstance.Trusted

	// upon event ⟨ uc, Propose | v ⟩ do

	//upon event ⟨ ec, StartEpoch | newts' , newl' ⟩ do

	//upon event ⟨ ep. ts , Aborted | state ⟩ such that ts = ets do

	//upon l = self ∧ val = ⊥ ∧ proposed = FALSE do

	//upon event ⟨ ep. ts , Decide | v ⟩ such that ts = ets do

	return uc
}
