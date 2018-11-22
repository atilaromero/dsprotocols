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
package leadership

import (
	"math/rand"
	"time"
)

/*
	Module:
		Name: EventualLeaderDetector, instance Ω.
	Events:
		Indication: (Ω, Trust | p): Indicates that process p is trusted to be leader.
	Properties:
		ELD1: Eventual accuracy: There is a time after which every correct process trusts
		some correct process.
		ELD2: Eventual agreement: There is a time after which no two correct processes
		trust different correct processes.
*/

// TrustMsg contains the ID of the current leader process.
type TrustMsg struct {
	ID int
}

// LeaderDetector (Eventual Leader Detector) is a struct that contains:
// Ind: deliver TrustMsg indicating the new leader.
type LeaderDetector struct {
	Suspected []int
	Leader    int
	Ind       chan TrustMsg
}

// NewLeaderDetector returns a NewLeaderDetector struct, which implements (fake) Eventual Leader Detector.
//
// There is 1 channel here: ind is used to indicate a leader to the upper layer.
// Receives the number of processes
func NewLeaderDetector(numProc int) LeaderDetector {

	susp := []int{}
	ind := make(chan TrustMsg)
	ld := LeaderDetector{susp, 0, ind}

	go func() {

		// TODO
		// events that need to be implemented based on indications from an eventual failure detector

		// when receiving an indication of suspect:
		//upon event  3P , Suspect | p  do
		//	suspected := suspected ∪ {p} ;

		// when receiving an indication of restore:
		//upon event  3P , Restore | p  do
		//	suspected := suspected \ {p} ;

		//first leader is the process id 0
		id := 0
		// deliver the trust message one layer up
		ld.Ind <- TrustMsg{
			ID: id,
		}
		for {
			rand.Seed(time.Now().Unix())
			random := rand.Intn(10)
			// 10% chances of changing the leader
			if random == 1 {
				// chooses another random leader
				newid := rand.Intn(numProc)
				if newid != id {
					id = newid
					// deliver the trust message one layer up
					ld.Ind <- TrustMsg{
						ID: id,
					}
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()

	return ld
}
