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

/*
	Algorithm 5.7: Leader-Driven Consensus
	Implements:
		UniformConsensus, instance uc.
	Uses:
		EpochChange, instance ec;
		EpochConsensus (multiple instances).
*/

type UcProposeMsg struct {
	Val int
}

type UcDeliverMsg struct {
	Val int
}

// Uc (Uniform Consensus,) is a struct that contains:
// EcInstance: epoch change instance
// map of ep: multiple instances of epoch consensus
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

func NewUC(ec *Ec, ep *Ep) *Uc {

	uc := &Uc{
		Req:        make(chan UcProposeMsg),
		Ind:        make(chan UcDeliverMsg),
		EcInstance: ec,
		EpInstance: ep,
		Val:        -1,
		Proposed:   false,
		Decided:    false,
		Ets:        -1,
		L:          -1,
		NewTS:      -1,
		NewL:       -1,
	}
	return uc
}

func (uc *Uc) Init() {
	// upon event ⟨ uc, Init ⟩ do
	// 		val := ⊥ ;
	// 		proposed := FALSE ; decided := FALSE ;
	// 		Obtain the leader l0 of the initial epoch with timestamp 0 from epoch-change inst. ec;
	// 		Initialize a new instance ep.0 of epoch consensus with timestamp 0,
	// 			leader l0 , and state (0, ⊥) ;
	// 		(ets, l) := (0, l0 ) ;
	// 		(newts, newl) := (0, ⊥) ;
	uc.Val = -1
	uc.Proposed = false
	uc.Decided = false
	l0 := uc.EcInstance.Trusted
	uc.EpInstance.Init(l0, State{ValTS: 0, Val: -1})
	uc.Ets = 0
	uc.L = l0
	uc.NewTS = 0
	uc.NewL = -1

	go func() {
		for {
			select {
			// upon event ⟨ uc, Propose | v ⟩ do
			// 		val := v ;
			case msg, ok := <-uc.Req:
				if !ok {
					return
				}
				uc.Val = msg.Val
				maybePropose(uc)

			// upon event ⟨ ec, StartEpoch | newts' , newl' ⟩ do
			// 		(newts, newl) := (newts' , newl') ;
			// 		trigger <ep. ets , Abort> ;
			case msg, ok := <-uc.EcInstance.Ind:
				if !ok {
					return
				}
				uc.NewTS = msg.Ts
				uc.NewL = msg.Leader
				uc.EpInstance.Req <- EpProposeMsg{Abort: true}

			case msg, ok := <-uc.EpInstance.Ind:
				if !ok {
					return
				}
				// upon event ⟨ ep. ts , Aborted | state ⟩ such that ts = ets do
				// 		(ets, l) := (newts, newl) ;
				// 		proposed := FALSE ;
				// 		Initialize a new instance ep. ets of epoch consensus with timestamp ets ,
				// 		leader l , and state state;
				if msg.Abort {
					uc.Ets = uc.NewTS
					uc.L = uc.NewL
					uc.Proposed = false
					uc.EpInstance.Init(uc.L, State{ValTS: uc.Ets, Val: msg.State.Val})
					continue
				}
				// upon event ⟨ ep. ts , Decide | v ⟩ such that ts = ets do
				// 		if decided = FALSE then
				// 			decided := TRUE ;
				// 			trigger < uc, Decide | v > ;
				if msg.State.ValTS != uc.Ets {
					continue
				}
				if uc.Decided {
					continue
				}
				uc.Decided = true
				uc.Ind <- UcDeliverMsg{msg.State.Val}
			}
		}
	}()
}

// upon l = self ∧ val != ⊥ ∧ proposed = FALSE do
// 		proposed := TRUE ;
// 		trigger < ep. ets , Propose | val > ;
func maybePropose(uc *Uc) {
	if uc.L == uc.EcInstance.Pl.ID() &&
		uc.Val != -1 &&
		!uc.Proposed {

		uc.Proposed = true
		uc.EpInstance.Req <- EpProposeMsg{Abort: true, Val: uc.Val}
		<-uc.EpInstance.Ind
		uc.EpInstance.Init(uc.L, State{ValTS: uc.Ets, Val: -1})
		uc.EpInstance.Req <- EpProposeMsg{Abort: false, Val: uc.Val}
	}
}
