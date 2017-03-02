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
