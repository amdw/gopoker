/*
Copyright 2013, 2015, 2017 Andrew Medworth

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
package holdem

import (
	"fmt"
	"github.com/amdw/gopoker/poker"
	"math"
	"testing"
)

func assertPotsWonSanity(winCount int, potsWon float64, description string, t *testing.T) {
	if potsWon > float64(winCount) || potsWon < 0 || (winCount > 0 && math.Abs(potsWon) < 1e-6) {
		t.Errorf("Illogical pot win total %v for %v (win count %v)", potsWon, description, winCount)
	}
}

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
	assertPotsWonSanity(sim.WinCount, sim.PotsWon, "us", t)
	assertPotsWonSanity(sim.BestOpponentWinCount, sim.BestOpponentPotsWon, "best opponent", t)
	assertPotsWonSanity(sim.RandomOpponentWinCount, sim.RandomOpponentPotsWon, "random opponent", t)
	if sim.PotsWon+sim.BestOpponentPotsWon > float64(simulations) {
		t.Errorf("More pots won than there were simulated hands: %v+%v vs %v", sim.PotsWon, sim.BestOpponentPotsWon, simulations)
	}
	betBreakEven := sim.PotOddsBreakEven()
	if betBreakEven < 0 || math.IsInf(betBreakEven, -1) || math.IsNaN(betBreakEven) {
		t.Errorf("Illogical pot odds break-even point: %v", betBreakEven)
	}
	checkCounts := func(counts []int, shouldSumToSims bool, name string) int {
		if len(counts) != int(poker.MAX_HANDCLASS) {
			t.Errorf("Expected %v %v, found %v", poker.MAX_HANDCLASS, name, len(counts))
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
		if poker.Beats(l, sim.BestHand) {
			t.Errorf("Best hand %v of class %v better than overall best %v", l, c, sim.BestHand)
		}
	}
	for c, l := range sim.ClassBestOppHands {
		if poker.Beats(l, sim.BestOppHand) {
			t.Errorf("Best opponent hand %v of class %v better than overall best %v", l, c, sim.BestOppHand)
		}
	}
	checkTiebreaks := func(tbs []poker.Rank, name string) {
		if len(tbs) != 5 {
			t.Errorf("Expected 5 tiebreaks for %v, found %v", name, len(tbs))
		}
	}
	// Catches error with best-hand zero value
	checkTiebreaks(sim.ClassBestHands[poker.HighCard].Tiebreaks, "high-card best hands")
	checkTiebreaks(sim.ClassBestOppHands[poker.HighCard].Tiebreaks, "high-card opponent best hands")
}

func TestSimSanity(t *testing.T) {
	sim := Simulator{}
	players := 5
	simulations := 10000
	sim.SimulateHoldem([]poker.Card{}, []poker.Card{}, players, simulations)
	assertSimSanity(&sim, players, simulations, t)
}

func TestTwoPlayers(t *testing.T) {
	sim := Simulator{}
	simulations := 10000
	sim.SimulateHoldem([]poker.Card{}, []poker.Card{}, 2, simulations)
	assertSimSanity(&sim, 2, simulations, t)
	// We can make some extra assertions here, as it's impossible for a pot to be split among opponents
	totalPotsWon := sim.PotsWon + sim.BestOpponentPotsWon
	if math.Abs(totalPotsWon-float64(simulations)) > 1e-6 {
		t.Errorf("Total pots won does not add up: %v us + %v best opponent = %v vs %v", sim.PotsWon, sim.BestOpponentPotsWon, totalPotsWon, simulations)
	}
	totalPotsWon = sim.PotsWon + sim.RandomOpponentPotsWon
	if math.Abs(totalPotsWon-float64(simulations)) > 1e-6 {
		t.Errorf("Total pots won does not add up: %v us + %v random opponent = %v vs %v", sim.PotsWon, sim.RandomOpponentPotsWon, totalPotsWon, simulations)
	}
}

func sp(r1s, r2s string, suited bool) StartingPair {
	r1, err := poker.MakeRank(r1s)
	if err != nil {
		panic(fmt.Sprintf("Cannot make rank from %v", r1s))
	}
	r2, err := poker.MakeRank(r2s)
	if err != nil {
		panic(fmt.Sprintf("Cannot make rank from %v", r2s))
	}
	return StartingPair{r1, r2, suited}
}

func TestPairs(t *testing.T) {
	pairs := []StartingPair{sp("K", "Q", false), sp("K", "Q", true), sp("K", "K", false)}
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

func TestEnumeration(t *testing.T) {
	sim := Simulator{}
	yourCards := h("9D", "7C")
	tableCards := h("KS", "7D", "AH", "8C", "8D")
	sim.SimulateHoldem(tableCards, yourCards, 2, 10000)

	// Should only do 45C2 = 990 simulations, one for each possible hand our opponent holds
	assertSimSanity(&sim, 2, 990, t)

	if sim.OurClassCounts[poker.TwoPair] != 990 {
		t.Errorf("We got two pair but got %v not 990!", sim.OurClassCounts[poker.TwoPair])
	}
	if sim.BestOpponentClassCounts[poker.FourOfAKind] != 1 {
		t.Errorf("Opponent has one way to make quads but found %v", sim.BestOpponentClassCounts[poker.FourOfAKind])
	}
	// There are two 8s, two 7s, three As and three Ks which could be in our opponent's hand.
	// So they have the following ways to make a full house: 3 AAs, 3 KKs, one 77, six A8s, six K8s, four 78s
	// for a total of 23.
	if sim.BestOpponentClassCounts[poker.FullHouse] != 23 {
		t.Errorf("Opponent has 23 ways to make a full house but found %v", sim.BestOpponentClassCounts[poker.FullHouse])
	}
	// There are two 8s left in the deck and 43 other cards, so 86 ways to get exactly three 8s.
	// Of these, 16 also pair another table card giving a full house, so there are 70 ways to be "on trips".
	if sim.BestOpponentClassCounts[poker.ThreeOfAKind] != 70 {
		t.Errorf("Opponent has 70 ways to make trips but found %v", sim.BestOpponentClassCounts[poker.ThreeOfAKind])
	}
	// 3 kings * 40 non-K/8 cards = 120 ways to get kings and another pair
	// 2 sevens * 38 non-7/8/K cards = 76 ways to get sevens and another (non-K) pair
	// 3 aces * 35 non-A/8/K/7 cards = 105 ways to get aces and another (non-K/7) pair
	// 4C2 = 6 pairs for each of 8 values not represented, plus 3C2 = 3 nine pairs, for a total of 51 pocket pairs
	// Total 352 ways to make two pair
	if sim.BestOpponentClassCounts[poker.TwoPair] != 352 {
		t.Errorf("Opponent has 352 ways to make two pair but found %v", sim.BestOpponentClassCounts[poker.TwoPair])
	}
	// And 990 minus all the other scenarios gives 544
	if sim.BestOpponentClassCounts[poker.OnePair] != 544 {
		t.Errorf("Opponent has 544 ways to make one pair but found %v", sim.BestOpponentClassCounts[poker.OnePair])
	}
}
