/*
Copyright 2013 Andrew Medworth

This file is part of Gopoker, a set of miscellaneous poker-related functions
written in the Go programming language (http://golang.org).

Gopoker is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Gopoker is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with Gopoker.  If not, see <http://www.gnu.org/licenses/>.
*/
package poker

import (
	"math"
	"testing"
)

func assertSimSanity(sim *Simulator, players, simulations int, t *testing.T) {
	if sim.Players != players {
		t.Errorf("Expected %v players, found %v", players, sim.Players)
	}
	if sim.HandCount != simulations {
		t.Errorf("Expected %v found %v for HandCount", simulations, sim.HandCount)
	}
	if sim.WinCount < 0 || sim.WinCount > simulations {
		t.Errorf("Illogical win count %v", sim.WinCount)
	}
	if sim.PotsWon > float64(sim.WinCount) || sim.PotsWon < 0 || (sim.WinCount > 0 && math.Abs(sim.PotsWon) < 1e-6) {
		t.Errorf("Illogical pot win total %v (win count %v)", sim.PotsWon, sim.WinCount)
	}
	betBreakEven := sim.PotOddsBreakEven()
	if betBreakEven < 0 || math.IsInf(betBreakEven, -1) || math.IsNaN(betBreakEven) {
		t.Errorf("Illogical pot odds break-even point: %v", betBreakEven)
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
	checkCounts(sim.BestOpponentClassCounts, true, "BestOpponentClassCounts")
	checkCounts(sim.RandomOpponentClassCounts, true, "RandomOpponentClassCounts")
	ourWins := checkCounts(sim.ClassWinCounts, false, "ClassWinCounts")
	jointWins := checkCounts(sim.ClassJointWinCounts, false, "ClassJointWinCounts")
	bestOppWins := checkCounts(sim.ClassBestOppWinCounts, false, "ClassBestOppWinCounts")
	if ourWins != sim.WinCount {
		t.Errorf("Class win counts should sum to %v, found %v", sim.WinCount, ourWins)
	}
	if jointWins != sim.JointWinCount {
		t.Errorf("Class joint win counts should sum to %v, found %v", sim.JointWinCount, jointWins)
	}
	if bestOppWins != sim.BestOpponentWinCount {
		t.Errorf("Best opponent win counts should sum to %v, found %v", sim.BestOpponentWinCount, bestOppWins)
	}
	if ourWins+bestOppWins-sim.JointWinCount != simulations {
		t.Errorf("Our wins (%v) and opponent wins (%v) minus joint wins (%v) sum to %v, expected %v", ourWins, bestOppWins, sim.JointWinCount, ourWins+bestOppWins-sim.JointWinCount, simulations)
	}
	randOppWins := checkCounts(sim.ClassRandOppWinCounts, false, "ClassRandOppWinCounts")
	if randOppWins != sim.RandomOpponentWinCount {
		t.Errorf("Random opponent wins %v but classes sum to %v", sim.RandomOpponentWinCount, randOppWins)
	}
	if randOppWins > bestOppWins {
		t.Errorf("Random opponent won more than best opponent (%v vs %v)", randOppWins, bestOppWins)
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

func TestSimSanity(t *testing.T) {
	sim := &Simulator{}
	players := 5
	simulations := 10000
	sim.SimulateHoldem([]Card{}, []Card{}, players, simulations)
	assertSimSanity(sim, players, simulations, t)
}

func TestPairs(t *testing.T) {
	pairs := []StartingPair{StartingPair{King, Queen, false}, StartingPair{King, Queen, true}, StartingPair{King, King, false}}
	players := 6
	simCount := 1000
	for _, pair := range pairs {
		sim := pair.RunSimulation(players, simCount)
		assertSimSanity(sim, players, simCount, t)
	}
}

func TestPotOdds(t *testing.T) {
	sim := Simulator{}
	sim.HandCount = 10000
	tests := map[float64]float64{0.0: 0.0, float64(sim.HandCount) / 2.0: 1.0, float64(sim.HandCount): math.Inf(1)}
	for potsWon, expected := range tests {
		sim.PotsWon = potsWon
		breakEven := sim.PotOddsBreakEven()
		if breakEven != expected {
			t.Errorf("Expected pot odds break even %v, found %v", expected, breakEven)
		}
	}
}
