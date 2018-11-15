package leadership

import (
	"fmt"
	"testing"
)

func TestLeaderDetector(t *testing.T) {
	// start a leader detector for 10 processess
	leader := NewLeaderDetector(10)

	// when receiving a indication from leader detector
	for msg, ok := <-leader.Ind; ok; msg, ok = <-leader.Ind {
		// just print the leader
		fmt.Println("New leader: ", msg.ID)
	}
}
