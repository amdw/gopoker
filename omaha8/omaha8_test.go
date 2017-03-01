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
	"github.com/amdw/gopoker/poker"
	"math/rand"
	"reflect"
	"testing"
)

func checkDuplicateCard(overlap map[poker.Card]int, card poker.Card, t *testing.T) {
	if overlap[card] > 0 {
		t.Errorf("Found duplicate card %v", card)
	}
	overlap[card]++
}

func TestDeal(t *testing.T) {
	pack := poker.NewPack()
	randGen := rand.New(rand.NewSource(1234))
	pack.Shuffle(randGen)

	playerCount := 6
	tableCards, playerCards := Deal(&pack, playerCount)
	if len(tableCards) != 5 {
		t.Errorf("Expected five table cards, found %v", len(tableCards))
	}
	overlap := make(map[poker.Card]int)
	for _, card := range tableCards {
		checkDuplicateCard(overlap, card, t)
	}
	if len(playerCards) != playerCount {
		t.Errorf("Expected %v sets of player cards, found %v", playerCount, len(playerCards))
	}
	for i, hand := range playerCards {
		if len(hand) != 4 {
			t.Errorf("Expected four cards for %vth player, found %v", i, len(hand))
		}
		for _, card := range hand {
			checkDuplicateCard(overlap, card, t)
		}
	}
}

type classifyTest struct {
	tableCards, holeCards               []poker.Card
	expectedHighLevel, expectedLowLevel poker.HandLevel
	expectedLowQualifies                bool
}

var h = poker.TestMakeHand
var hl = poker.TestMakeHandLevel
var board = h("2S", "5C", "10H", "7D", "8C")
var classifyTests = []classifyTest{
	{board, h("AS", "4S", "5H", "KC"), hl("OnePair", "5", "A", "10", "8"), hl("HighCard", "7", "5", "4", "2", "A"), true},
	{board, h("AH", "3H", "10S", "10C"), hl("ThreeOfAKind", "10", "8", "7"), hl("HighCard", "7", "5", "3", "2", "A"), true},
	{board, h("7C", "9C", "JS", "QS"), hl("Straight", "J"), hl("HighCard", "9", "8", "7", "5", "2"), false},
	{board, h("4H", "6H", "KS", "KD"), hl("Straight", "8"), hl("HighCard", "7", "6", "5", "4", "2"), true},
	{board, h("AD", "3D", "6D", "9H"), hl("Straight", "10"), hl("HighCard", "7", "5", "3", "2", "A"), true},
}

func TestClassify(t *testing.T) {
	for _, test := range classifyTests {
		level := classify(test.tableCards, test.holeCards)
		if !reflect.DeepEqual(level.HighLevel, test.expectedHighLevel) {
			t.Errorf("Expected high level %q with board %q and hole cards %q, found %q", test.expectedHighLevel, test.tableCards, test.holeCards, level.HighLevel)
		}
		if !reflect.DeepEqual(level.LowLevel, test.expectedLowLevel) {
			t.Errorf("Expected low level %q with board %q and hole cards %q, found %q", test.expectedLowLevel, test.tableCards, test.holeCards, level.LowLevel)
		}
		if level.LowLevelQualifies != test.expectedLowQualifies {
			t.Errorf("Expected low level qualifies = %v with board %q and hole cards %q, found %v (%q)", test.expectedLowQualifies, test.tableCards, test.holeCards, level.LowLevelQualifies, level.LowLevel)
		}

		hl := poker.ClassifyHand(level.HighHand)
		if !reflect.DeepEqual(test.expectedHighLevel, hl) {
			t.Errorf("Expected high hand to classify as %q, found %q (%q)", test.expectedHighLevel, hl, level.HighHand)
		}
		ll := poker.ClassifyAceToFiveLow(level.LowHand)
		if !reflect.DeepEqual(test.expectedLowLevel, ll) {
			t.Errorf("Expected low hand to classify as %q, found %q (%q)", test.expectedLowLevel, ll, level.LowHand)
		}
	}
}
