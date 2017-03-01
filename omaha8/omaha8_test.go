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
	tableCards, holeCards []poker.Card
	expectedLevel         Omaha8Level
}

func o8l(highLevel, lowLevel poker.HandLevel) Omaha8Level {
	return Omaha8Level{highLevel, lowLevel, true}
}

func o8lh(highLevel poker.HandLevel) Omaha8Level {
	result := Omaha8Level{}
	result.HighLevel = highLevel
	result.LowLevelQualifies = false
	return result
}

var h = poker.TestMakeHand
var hl = poker.TestMakeHandLevel
var board = h("2S", "5C", "10H", "7D", "8C")
var classifyTests = []classifyTest{
	{board, h("AS", "4S", "5H", "KC"), o8l(hl("OnePair", "5", "A", "10", "8"), hl("HighCard", "7", "5", "4", "2", "A"))},
	{board, h("AH", "3H", "10S", "10C"), o8l(hl("ThreeOfAKind", "10", "8", "7"), hl("HighCard", "7", "5", "3", "2", "A"))},
	{board, h("7C", "9C", "JS", "QS"), o8lh(hl("Straight", "J"))},
	{board, h("4H", "6H", "KS", "KD"), o8l(hl("Straight", "8"), hl("HighCard", "7", "6", "5", "4", "2"))},
	{board, h("AD", "3D", "6D", "9H"), o8l(hl("Straight", "10"), hl("HighCard", "7", "5", "3", "2", "A"))},
}

func levelsEqual(l1, l2 Omaha8Level) bool {
	if reflect.DeepEqual(l1, l2) {
		return true
	}
	if !l1.LowLevelQualifies && !l2.LowLevelQualifies {
		return reflect.DeepEqual(l1.HighLevel, l2.HighLevel)
	}
	return false
}

func TestClassify(t *testing.T) {
	for _, test := range classifyTests {
		level := classify(test.tableCards, test.holeCards)
		if !levelsEqual(level, test.expectedLevel) {
			t.Errorf("Expected level %q with board %q and hole cards %q, found %q", test.expectedLevel, test.tableCards, test.holeCards, level)
		}
	}
}
