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
package omaha8

import (
	"fmt"
	"github.com/amdw/gopoker/poker"
	"math"
	"math/rand"
	"testing"
)

func TestFixedShuffle(t *testing.T) {
	tests := 1000
	players := 3
	pack := poker.NewPack()
	randGen := rand.New(rand.NewSource(1234))

	tableCards := h("AD", "QC", "6S", "QH", "3H")
	yourCards := h("3S", "4C", "5D", "6H")

	for tcPrefixLen := 3; tcPrefixLen <= 5; tcPrefixLen++ {
		tcPrefix := tableCards[0:tcPrefixLen]

		// Very crude randomness check
		randomHands := make(map[string]int)
		trackingLimit := 10

		for i := 0; i < tests; i++ {
			shuffleFixing(&pack, tcPrefix, yourCards, randGen)
			poker.TestPackPermutation(&pack, t)
			dealtTable, dealtPlayers := Deal(&pack, players)
			for j := 0; j < len(tcPrefix); j++ {
				if dealtTable[j] != tcPrefix[j] {
					t.Errorf("Expected %v at position %v, found %v", tcPrefix[j], j, dealtTable[j])
				}
			}
			if !poker.CardsEqual(yourCards, dealtPlayers[0]) {
				t.Errorf("Expected cards %v dealt to player 0, found %v", yourCards, dealtPlayers[0])
			}
			if len(randomHands) < trackingLimit {
				randomHands[fmt.Sprintf("%q", dealtPlayers[1])]++
			}
		}

		if len(randomHands) < trackingLimit {
			t.Errorf("Suspicious lack of randomness - only following hands seen: %q", randomHands)
		}
	}
}

func assertSimSanity(sim *Omaha8Simulator, players, simCount int, t *testing.T) {
	poker.TestAssertSimSanity(&sim.HighSimulator, players, simCount, t)
	totalPotsWon := sim.PotsWon()
	poker.TestAssertPotsWonSanity(sim.HighSimulator.WinCount+sim.LowSimulator.WinCount, totalPotsWon, "us (total)", t)
}

func TestSimulate(t *testing.T) {
	simCount := 10000
	players := 4
	yourCards := h("AS", "QC", "3D", "4H")
	randGen := rand.New(rand.NewSource(1234))
	sim := SimulateOmaha8([]poker.Card{}, yourCards, players, simCount, randGen)
	assertSimSanity(sim, players, simCount, t)
}

func TestPotOdds(t *testing.T) {
	sim := Omaha8Simulator{}
	sim.reset(2, 2)
	randGen := rand.New(rand.NewSource(1234))
	highWins := Omaha8Level{hl("StraightFlush", "A"), hl("HighCard", "8", "7", "6", "5", "4"),
		h("AS", "KS", "QS", "JS", "10S"), h("8C", "7D", "6S", "5H", "4C"), true}
	lowWins := Omaha8Level{hl("TwoPair", "A", "K", "Q"), hl("HighCard", "6", "4", "3", "2", "A"),
		h("AS", "AC", "KS", "KC", "QH"), h("6D", "4D", "3D", "2C", "AD"), true}
	outcome1 := PlayerOutcome{1, highWins, true, false, 0.5, 0.0}
	outcome2 := PlayerOutcome{2, lowWins, false, true, 0.0, 0.5}
	sim.processHand([]PlayerOutcome{outcome1, outcome2}, randGen)
	outcome1.Player = 2
	outcome2.Player = 1
	sim.processHand([]PlayerOutcome{outcome2, outcome1}, randGen)

	breakEven := sim.PotOddsBreakEven()
	if math.Abs(breakEven-1.0) > 1e-6 {
		t.Errorf("Expected even pot odds, found %v", breakEven)
	}
}
