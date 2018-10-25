package consensus

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
	EcInstance Ec
	Consensus  map[int]Ep
}

func NewUC() *Uc {

	//upon event ⟨ uc, Init ⟩ do
	uc := Uc{
		NewEc(nil, nil, nil, 0),
		make(map[int]Ep),
	}

	// upon event ⟨ uc, Propose | v ⟩ do

	//upon event ⟨ ec, StartEpoch | newts' , newl' ⟩ do

	//upon event ⟨ ep. ts , Aborted | state ⟩ such that ts = ets do

	//upon l = self ∧ val = ⊥ ∧ proposed = FALSE do

	//upon event ⟨ ep. ts , Decide | v ⟩ such that ts = ets do

	return &uc
}
