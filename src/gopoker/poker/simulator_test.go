package poker

import (
	"testing"
)

func TestSimSanity(t *testing.T) {
	sim := Simulator{}
	simulations := 10000
	sim.SimulateHoldem([]Card{}, []Card{}, 5, simulations)
	if sim.HandCount != simulations {
		t.Errorf("Expected %v found %v for HandCount", simulations, sim.HandCount)
	}
	if sim.WinCount < 0 || sim.WinCount > simulations {
		t.Errorf("Illogical win count %v", sim.WinCount)
	}
	checkCounts := func(counts []int, shouldSumToSims bool, name string) int {
		if len(counts) != int(MAX_HANDCLASS) {
			t.Errorf("Expected %v %v, found %v", MAX_HANDCLASS, name, len(counts))
		}
		sum := 0
		for i, c := range counts {
			if c < 0 || c > simulations {
				t.Errorf("Insane value %v at %v of %v", c, i, name)
			}
			sum += c
		}
		if sum > simulations {
			t.Errorf("Insane sum %v for %v", sum, name)
		}
		if shouldSumToSims && sum != simulations {
			t.Errorf("Expected sum %v for %v, found %v", simulations, name, sum)
		}
		return sum
	}
	checkCounts(sim.OurClassCounts, true, "OurClassCounts")
	checkCounts(sim.OpponentClassCounts, true, "OpponentClassCounts")
	ourWins := checkCounts(sim.ClassWinCounts, false, "ClassWinCounts")
	oppWins := checkCounts(sim.ClassOppWinCounts, false, "ClassOppWinCounts")
	if ourWins+oppWins != simulations {
		t.Errorf("Our wins and opponent wins sum to %v, expected %v", ourWins+oppWins, simulations)
	}
}
