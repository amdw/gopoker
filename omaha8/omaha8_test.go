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
	"math"
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
var hs = poker.TestMakeHands
var hls = poker.TestMakeHandLevels
var board = h("2S", "5C", "10H", "7D", "8C")
var classifyTests = []classifyTest{
	{board, h("AS", "4S", "5H", "KC"), hl("OnePair", "5", "A", "10", "8"), hl("HighCard", "7", "5", "4", "2", "A"), true},
	{board, h("AH", "3H", "10S", "10C"), hl("ThreeOfAKind", "10", "8", "7"), hl("HighCard", "7", "5", "3", "2", "A"), true},
	{board, h("7C", "9C", "JS", "QS"), hl("Straight", "J"), hl("HighCard", "9", "8", "7", "5", "2"), false},
	{board, h("4H", "6H", "KS", "KD"), hl("Straight", "8"), hl("HighCard", "7", "6", "5", "4", "2"), true},
	{board, h("AD", "3D", "6D", "9H"), hl("Straight", "10"), hl("HighCard", "7", "5", "3", "2", "A"), true},
	{h("6S", "7S", "8C", "JD", "QH"), h("AS", "3S", "KS", "KC"), hl("OnePair", "K", "Q", "J", "8"), hl("HighCard", "8", "7", "6", "3", "A"), true},
	{h("AS", "2S", "3S", "5S", "5C"), h("5D", "5H", "6C", "7D"), hl("FourOfAKind", "5", "A"), hl("HighCard", "6", "5", "3", "2", "A"), true},
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

type outcomeTest struct {
	tableCards          []poker.Card
	playerCards         [][]poker.Card
	expectedHighs       []poker.HandLevel
	expectedLows        []poker.HandLevel
	expectedHighWinners []bool
	expectedLowWinners  []bool
	expectedPotSplit    []float64
}

func split(fracs ...float64) []float64 {
	return fracs
}

func bs(bools ...bool) []bool {
	return bools
}

var outcomeTests = []outcomeTest{
	{h("6S", "7S", "8C", "JD", "QH"), hs(h("AS", "3S", "KS", "KC"), h("2S", "3C", "5D", "9H")),
		hls(hl("OnePair", "K", "Q", "J", "8"), hl("Straight", "9")),
		hls(hl("HighCard", "8", "7", "6", "3", "A"), hl("HighCard", "8", "7", "6", "3", "2")),
		bs(false, true), bs(true, false), split(0.5, 0.5)},
	{h("2S", "3C", "4D", "5H", "JS"), hs(h("AS", "6C", "KD", "KH"), h("AC", "4S", "8D", "8H")),
		hls(hl("OnePair", "K", "J", "5", "4"), hl("Straight", "5")),
		hls(hl("HighCard", "6", "4", "3", "2", "A"), hl("HighCard", "5", "4", "3", "2", "A")),
		bs(false, true), bs(false, true), split(0.0, 1.0)},
	{h("AS", "2C", "3D", "4H", "5S"), hs(h("AC", "3S", "5D", "6S"), h("3C", "4D", "JS", "JC")),
		hls(hl("Straight", "6"), hl("Straight", "5")),
		hls(hl("HighCard", "5", "4", "3", "2", "A"), hl("HighCard", "5", "4", "3", "2", "A")),
		bs(true, false), bs(true, true), split(0.75, 0.25)},
	{h("2S", "4C", "6D", "7H", "8S"), hs(h("AS", "AC", "4D", "7D"), h("2C", "5S", "5C", "7S")),
		hls(hl("TwoPair", "7", "4", "8"), hl("Straight", "8")),
		hls(hl("HighCard", "7", "6", "4", "2", "A"), hl("HighCard", "7", "6", "5", "4", "2")),
		bs(false, true), bs(true, false), split(0.5, 0.5)},
	{h("2S", "4C", "7D", "8H", "9S"), hs(h("AS", "AC", "4D", "7H"), h("2C", "5S", "6C", "KD")),
		hls(hl("TwoPair", "7", "4", "9"), hl("Straight", "9")),
		hls(hl("HighCard", "8", "7", "4", "2", "A"), hl("HighCard", "7", "6", "5", "4", "2")),
		bs(false, true), bs(false, true), split(0.0, 1.0)},
	{h("3S", "9C", "10D", "JH", "QS"), hs(h("AS", "3C", "JC", "JD"), h("4S", "4C", "8S", "QH")),
		hls(hl("ThreeOfAKind", "J", "Q", "10"), hl("Straight", "Q")),
		hls(hl("HighCard", "J", "10", "9", "3", "A"), hl("HighCard", "10", "9", "8", "4", "3")), // Non-qualifying
		bs(false, true), bs(false, false), split(0.0, 1.0)},
	{h("6H", "7H", "9H", "JH", "KD"), hs(h("AH", "2C", "8C", "10D"), h("3S", "4C", "5H", "8D")),
		hls(hl("Straight", "J"), hl("Straight", "9")),
		hls(hl("HighCard", "9", "7", "6", "2", "A"), hl("HighCard", "9", "7", "6", "4", "3")), // Non-qualifying
		bs(true, false), bs(false, false), split(1.0, 0.0)},
	{h("3S", "3C", "7D", "JH", "JS"), hs(h("AS", "2C", "4D", "JH"), h("4S", "5C", "7C", "7H")),
		hls(hl("ThreeOfAKind", "J", "A", "7"), hl("FullHouse", "7", "J")),
		hls(hl("HighCard", "J", "7", "3", "2", "A"), hl("HighCard", "J", "7", "5", "4", "3")), // Non-qualifying
		bs(false, true), bs(false, false), split(0.0, 1.0)},
	{h("3S", "3C", "7D", "JH", "JS"), hs(h("AS", "2C", "3D", "JH"), h("4S", "5C", "7C", "7H")),
		hls(hl("FullHouse", "J", "3"), hl("FullHouse", "7", "J")),
		hls(hl("HighCard", "J", "7", "3", "2", "A"), hl("HighCard", "J", "7", "5", "4", "3")), // Non-qualifying
		bs(true, false), bs(false, false), split(1.0, 0.0)},
	{h("8S", "8C", "8D", "8H", "9S"), hs(h("2S", "3C", "3D", "QH"), h("AS", "2C", "3H", "QC")),
		hls(hl("FullHouse", "8", "3"), hl("ThreeOfAKind", "8", "A", "Q")),
		hls(hl("OnePair", "8", "9", "3", "2"), hl("OnePair", "8", "9", "2", "A")), // Non-qualifying
		bs(true, false), bs(false, false), split(1.0, 0.0)},
	{h("AS", "4C", "5D", "8H", "9S"), hs(h("AC", "4S", "5H", "8D"), h("AD", "4D", "5C", "KS")),
		hls(hl("TwoPair", "A", "8", "9"), hl("TwoPair", "A", "5", "9")),
		hls(hl("HighCard", "9", "8", "5", "4", "A"), hl("HighCard", "9", "8", "5", "4", "A")), // Non-qualifying
		bs(true, false), bs(false, false), split(1.0, 0.0)},
	{h("AS", "4C", "5D", "8H", "2S"), hs(h("AC", "4S", "5H", "8D"), h("AD", "4D", "5C", "KS")),
		hls(hl("TwoPair", "A", "8", "5"), hl("TwoPair", "A", "5", "8")),
		hls(hl("HighCard", "8", "5", "4", "2", "A"), hl("HighCard", "8", "5", "4", "2", "A")), // Now qualifying equal
		bs(true, false), bs(true, true), split(0.75, 0.25)},
}

func TestPlayerOutcomes(t *testing.T) {
	for _, test := range outcomeTests {
		outcomes := PlayerOutcomes(test.tableCards, test.playerCards)
		for i, outcome := range outcomes {
			if !reflect.DeepEqual(test.expectedHighs[i], outcome.Level.HighLevel) {
				t.Errorf("Expected high %q for player @ %v, found %q on board %q hands %q", test.expectedHighs[i], i, outcome.Level.HighLevel, test.tableCards, test.playerCards)
			}
			if !reflect.DeepEqual(test.expectedLows[i], outcome.Level.LowLevel) {
				t.Errorf("Expected high %q for player @ %v, found %q on board %q hands %q", test.expectedLows[i], i, outcome.Level.LowLevel, test.tableCards, test.playerCards)
			}
			if test.expectedHighWinners[i] != outcome.IsHighWinner {
				t.Errorf("Expected high winner = %v for player @ %v, found %v on board %q hands %q", test.expectedHighWinners[i], i, outcome.IsHighWinner, test.tableCards, test.playerCards)
			}
			if test.expectedLowWinners[i] != outcome.IsLowWinner {
				t.Errorf("Expected low winner = %v for player @ %v, found %v on board %q hands %q", test.expectedLowWinners[i], i, outcome.IsLowWinner, test.tableCards, test.playerCards)
			}
			if math.Abs(test.expectedPotSplit[i]-outcome.PotFractionWon()) > 1e-6 {
				t.Errorf("Expected pot split %v for player @ %v on board %q hand %q, found %v", test.expectedPotSplit[i], i, test.tableCards, test.playerCards, outcome.PotFractionWon())
			}
		}
	}
}
