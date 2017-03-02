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
	"github.com/amdw/gopoker/poker"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

var h = poker.TestMakeHand
var hs = poker.TestMakeHands
var hl = poker.TestMakeHandLevel

func levelsEqual(l1, l2 poker.HandLevel) bool {
	return reflect.DeepEqual(l1, l2)
}

type classificationTest struct {
	HandCards     []poker.Card
	TableCards    []poker.Card
	ExpectedLevel poker.HandLevel
	ExpectedCards []poker.Card
}

var classificationTests = []classificationTest{
	{h("10S", "JS"), h("2H", "QS", "6D", "KS", "AS"), hl("StraightFlush", "A"), h("10S", "JS", "QS", "KS", "AS")},
	{h("AH", "3H"), h("2H", "4H", "5H", "2C", "3C"), hl("StraightFlush", "5"), h("AH", "2H", "3H", "4H", "5H")},
	{h("10S", "JS"), h("10C", "QD", "10D", "10H", "2S"), hl("FourOfAKind", "10", "Q"), h("10S", "10C", "10D", "10H", "QD")},
	{h("10S", "10C"), h("10D", "10H", "JD", "QD", "KD"), hl("FourOfAKind", "10", "K"), h("10S", "10C", "10D", "10H", "KD")},
	{h("5C", "10S"), h("5S", "5D", "5H", "AS", "KC"), hl("FourOfAKind", "5", "A"), h("5C", "5S", "5D", "5H", "AS")},
	{h("2S", "2H"), h("3D", "3H", "3C", "QD", "KS"), hl("FullHouse", "3", "2"), h("2S", "2H", "3H", "3C", "3D")},
	{h("2S", "3S"), h("4C", "4S", "2D", "2H", "3C"), hl("FullHouse", "2", "4"), h("2S", "2D", "2H", "4C", "4S")},
	{h("6H", "8H"), h("9H", "10H", "2H", "3S", "7C"), hl("Flush", "10", "9", "8", "6", "2"), h("10H", "9H", "8H", "6H", "2H")},
	{h("6S", "8H"), h("9H", "10H", "JH", "QH", "7H"), hl("StraightFlush", "Q"), h("QH", "JH", "10H", "9H", "8H")},
	{h("AS", "JH"), h("QC", "KD", "10S", "2C", "3C"), hl("Straight", "A"), h("AS", "KD", "QC", "JH", "10S")},
	{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), hl("Straight", "5"), h("AS", "2C", "3H", "4C", "5D")},
	{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), hl("ThreeOfAKind", "6", "K", "J"), h("6S", "6D", "6C", "KH", "JC")},
	{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), hl("ThreeOfAKind", "6", "K", "J"), h("6S", "6D", "6C", "KH", "JC")},
	{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), hl("TwoPair", "6", "4", "A"), h("6S", "6D", "4D", "4S", "AH")},
	{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), hl("TwoPair", "6", "4", "A"), h("6S", "6D", "4D", "4S", "AH")},
	{h("AS", "AH"), h("2S", "4C", "6D", "8S", "10D"), hl("OnePair", "A", "10", "8", "6"), h("AS", "AH", "10D", "8S", "6D")},
	{h("AS", "2S"), h("AH", "4C", "6D", "8S", "10D"), hl("OnePair", "A", "10", "8", "6"), h("AS", "AH", "10D", "8S", "6D")},
	{h("2S", "4S"), h("5D", "7S", "8S", "QH", "KH"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	{h("2S", "KH"), h("5D", "7S", "8S", "QH", "4S"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	{h("8S", "KH"), h("5D", "7S", "2S", "QH", "4S"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	{h("KS", "8C"), h("QH", "10C", "9H", "7H", "6S"), hl("Straight", "10"), h("10C", "9H", "8C", "7H", "6S")},
	{h("KS", "10H"), h("8C", "QC", "9H", "7H", "6S"), hl("Straight", "10"), h("10H", "9H", "8C", "7H", "6S")},
	{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), hl("Straight", "10"), h("10H", "9H", "8C", "7H", "6S")},
	// Uses none of the hand cards
	{h("2S", "3S"), h("AH", "10H", "JH", "QH", "KH"), hl("StraightFlush", "A"), h("AH", "KH", "QH", "JH", "10H")},
	// Uses one of the hand cards
	{h("AH", "3S"), h("2S", "10H", "JH", "QH", "KH"), hl("StraightFlush", "A"), h("AH", "KH", "QH", "JH", "10H")},
}

func TestClassification(t *testing.T) {
	for i, test := range classificationTests {
		level, cards := classify(test.TableCards, test.HandCards)
		if !reflect.DeepEqual(test.ExpectedLevel, level) {
			t.Errorf("Test @%v: expected level %v, found %v (hand %q table %q)", i, test.ExpectedLevel, level, test.HandCards, test.TableCards)
		}
		if !poker.CardsEqual(test.ExpectedCards, cards) {
			t.Errorf("Test @%v: expected cards %v, found %v (hand %q table %q)", i, test.ExpectedCards, cards, test.HandCards, test.TableCards)
		}
	}
}

func cardDupeCheck(cards []poker.Card, counts map[poker.Card]int, t *testing.T) {
	for _, card := range cards {
		counts[card]++
		if counts[card] > 1 {
			t.Errorf("Found duplicate card %v", card)
		}
	}
}

func TestDeal(t *testing.T) {
	p := poker.NewPack()
	randGen := rand.New(rand.NewSource(1234))
	p.Shuffle(randGen)

	players := 6
	tableCards, playerCards := Deal(&p, players)

	if len(tableCards) != 5 {
		t.Errorf("Expected 5 table cards, found %v", len(tableCards))
	}
	if len(playerCards) != players {
		t.Errorf("Expected %v player hands, found %v", len(playerCards))
	}
	dupeCheck := make(map[poker.Card]int)
	cardDupeCheck(tableCards, dupeCheck, t)
	for _, hand := range playerCards {
		if len(hand) != 2 {
			t.Errorf("Expected 2 cards in player hand, found %v", len(hand))
		}
		cardDupeCheck(hand, dupeCheck, t)
	}
}

func po(level poker.HandLevel, cards []poker.Card, won bool, potFraction float64) PlayerOutcome {
	return PlayerOutcome{0, level, cards, won, potFraction}
}

func pos(outcomes ...PlayerOutcome) []PlayerOutcome {
	result := make([]PlayerOutcome, len(outcomes))
	for i, outcome := range outcomes {
		result[i] = outcome
		result[i].Player = i + 1
	}
	return result
}

type gameOutcomeTest struct {
	tableCards       []poker.Card
	playerCards      [][]poker.Card
	expectedOutcomes []PlayerOutcome
}

var royalFlush = h("AS", "KS", "QS", "JS", "10S")
var royalFlushLevel = hl("StraightFlush", "A")
var gameOutcomeTests = []gameOutcomeTest{
	// Simple split case where both players play the board
	{royalFlush, hs(h("2C", "3C"), h("4C", "5C")),
		pos(po(royalFlushLevel, royalFlush, true, 0.5), po(royalFlushLevel, royalFlush, true, 0.5))},
	{royalFlush, hs(h("2C", "3C"), h("4C", "5C"), h("6C", "7C")),
		pos(po(royalFlushLevel, royalFlush, true, 1.0/3), po(royalFlushLevel, royalFlush, true, 1.0/3), po(royalFlushLevel, royalFlush, true, 1.0/3))},
	{h("2H", "3H", "4H", "5C", "JC"), hs(h("AH", "AC"), h("AD", "AS"), h("QD", "KD")),
		pos(po(hl("Straight", "5"), h("AH", "2H", "3H", "4H", "5C"), true, 0.5),
			po(hl("Straight", "5"), h("AD", "2H", "3H", "4H", "5C"), true, 0.5),
			po(hl("HighCard", "K", "Q", "J", "5", "4"), h("KD", "QD", "JC", "5C", "4H"), false, 0))},
	{h("2H", "3H", "4H", "5H", "JC"), hs(h("AH", "AC"), h("AD", "AS"), h("QD", "KD")),
		pos(po(hl("StraightFlush", "5"), h("AH", "2H", "3H", "4H", "5H"), true, 1.0),
			po(hl("Straight", "5"), h("AD", "2H", "3H", "4H", "5H"), false, 0),
			po(hl("HighCard", "K", "Q", "J", "5", "4"), h("KD", "QD", "JC", "5H", "4H"), false, 0))},
}

func TestDealOutcomes(t *testing.T) {
	for _, test := range gameOutcomeTests {
		outcomes := DealOutcomes(test.tableCards, test.playerCards)
		if len(outcomes) != len(test.playerCards) {
			t.Errorf("Expected %v outcomes, found %v", len(test.playerCards), len(outcomes))
		}
		for i, outcome := range outcomes {
			if outcome.Player != i+1 {
				t.Errorf("Expected player %v at position %v, found %v", i+1, i, outcome.Player)
			}
			expectedOutcome := test.expectedOutcomes[i]
			if !reflect.DeepEqual(expectedOutcome.Level, outcome.Level) {
				t.Errorf("Expected level %q at position %v, found %q", expectedOutcome.Level, i, outcome.Level)
			}
			if !poker.CardsEqual(expectedOutcome.Cards, outcome.Cards) {
				t.Errorf("Expected cards %q at position %v, found %q", expectedOutcome.Cards, i, outcome.Cards)
			}
			if expectedOutcome.Won != outcome.Won {
				t.Errorf("Expected won = %v at position %v, found %v", expectedOutcome.Won, i, outcome.Won)
			}
			if math.Abs(expectedOutcome.PotFractionWon-outcome.PotFractionWon) > 1e-6 {
				t.Errorf("Expected pot fraction %v at position %v, found %v", expectedOutcome.PotFractionWon, i, outcome.PotFractionWon)
			}
			if outcome.Won && math.Abs(outcome.PotFractionWon) < 1e-6 {
				t.Errorf("Won at %v but didn't get any of the pot: %v", i, outcome.PotFractionWon)
			}
			if !outcome.Won && math.Abs(outcome.PotFractionWon) > 1e-6 {
				t.Errorf("Didn't win at %v but won %v of pot", i, outcome.PotFractionWon)
			}
		}

		potSum := 0.0
		for _, outcome := range outcomes {
			potSum += outcome.PotFractionWon
		}
		if math.Abs(potSum-1.0) > 1e-6 {
			t.Errorf("Pot fractions add up to %v, expected 1.0", potSum)
		}
	}
}
