/*
Copyright 2017 Andrew Medworth

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
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestMakeHand(cards ...string) []Card {
	result := make([]Card, len(cards))
	for i, c := range cards {
		result[i] = C(c)
	}
	return result
}

var h = TestMakeHand // Helper for tests in this package

func TestMakeHands(hands ...[]Card) [][]Card {
	return hands
}

func parseHandClass(handClassStr string) HandClass {
	switch handClassStr {
	case "StraightFlush":
		return StraightFlush
	case "FourOfAKind":
		return FourOfAKind
	case "FullHouse":
		return FullHouse
	case "Flush":
		return Flush
	case "Straight":
		return Straight
	case "ThreeOfAKind":
		return ThreeOfAKind
	case "TwoPair":
		return TwoPair
	case "OnePair":
		return OnePair
	case "HighCard":
		return HighCard
	default:
		panic(fmt.Sprintf("Unknown hand class %v", handClassStr))
	}
}

func TestMakeHandLevel(handClassStr string, tieBreakRankStrs ...string) HandLevel {
	class := parseHandClass(handClassStr)
	tieBreaks := make([]Rank, len(tieBreakRankStrs))
	for i, rankStr := range tieBreakRankStrs {
		var err error
		tieBreaks[i], err = MakeRank(rankStr)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse rank %v", rankStr))
		}
	}
	return HandLevel{class, tieBreaks}
}

var hl = TestMakeHandLevel // Helper for tests in this package

func TestMakeHandLevels(levels ...HandLevel) []HandLevel {
	return levels
}

// Check equality of sets of cards, ignoring ordering (mutates inputs)
func CardsEqual(c1, c2 []Card) bool {
	SortCards(c1, false)
	SortCards(c2, false)
	return reflect.DeepEqual(c1, c2)
}

// Assert that the pack contains exactly one of every card
func TestPackPermutation(pack *Pack, t *testing.T) {
	permCheck := make([][]int, 4)
	for i := 0; i < 4; i++ {
		permCheck[i] = make([]int, 13)
	}
	for _, c := range pack.Cards {
		permCheck[c.Suit][c.Rank]++
	}
	for s := range permCheck {
		for r, count := range permCheck[s] {
			if count != 1 {
				t.Fatalf("Expected exactly one %v%v in pack after shuffle, found %v", Rank(r).String(), Suit(s).String(), count)
			}
		}
	}

}
func TestAssertPotsWonSanity(winCount int, potsWon float64, description string, t *testing.T) {
	if potsWon > float64(winCount) || potsWon < 0 || (winCount > 0 && math.Abs(potsWon) < 1e-6) {
		t.Errorf("Illogical pot win total %v for %v (win count %v)", potsWon, description, winCount)
	}
}

func TestAssertSimSanity(sim *Simulator, players, simulations int, t *testing.T) {
	if sim.Players != players {
		t.Errorf("Expected %v players, found %v", players, sim.Players)
	}
	if sim.HandCount != simulations {
		t.Errorf("Expected %v found %v for HandCount", simulations, sim.HandCount)
	}
	if sim.WinCount < 0 || sim.WinCount > simulations {
		t.Errorf("Illogical win count %v", sim.WinCount)
	}
	TestAssertPotsWonSanity(sim.WinCount, sim.PotsWon, "us", t)
	TestAssertPotsWonSanity(sim.BestOpponentWinCount, sim.BestOpponentPotsWon, "best opponent", t)
	TestAssertPotsWonSanity(sim.RandomOpponentWinCount, sim.RandomOpponentPotsWon, "random opponent", t)
	if sim.PotsWon+sim.BestOpponentPotsWon > float64(simulations) {
		t.Errorf("More pots won than there were simulated hands: %v+%v vs %v", sim.PotsWon, sim.BestOpponentPotsWon, simulations)
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
