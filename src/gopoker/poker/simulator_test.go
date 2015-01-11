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

	for c, l := range sim.ClassBestHands {
		if Beats(l, sim.BestHand) {
			t.Errorf("Best hand %v of class %v better than overall best %v", l, c, sim.BestHand)
		}
	}
	for c, l := range sim.ClassBestOppHands {
		if Beats(l, sim.BestOppHand) {
			t.Errorf("Best opponent hand %v of class %v better than overall best %v", l, c, sim.BestOppHand)
		}
	}
	checkTiebreaks := func(tbs []Rank, name string) {
		if len(tbs) != 5 {
			t.Errorf("Expected 5 tiebreaks for %v, found %v", name, len(tbs))
		}
	}
	// Catches error with best-hand zero value
	checkTiebreaks(sim.ClassBestHands[HighCard].Tiebreaks, "high-card best hands")
	checkTiebreaks(sim.ClassBestOppHands[HighCard].Tiebreaks, "high-card opponent best hands")
}