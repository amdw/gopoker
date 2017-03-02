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
	"math/rand"
	"reflect"
	"testing"
)

var h = poker.TestMakeHand
var hl = poker.TestMakeHandLevel

func levelsEqual(l1, l2 poker.HandLevel) bool {
	return reflect.DeepEqual(l1, l2)
}

func cardsEqual(c1, c2 []poker.Card) bool {
	poker.SortCards(c1, false)
	poker.SortCards(c2, false)
	return reflect.DeepEqual(c1, c2)
}

type GameTest struct {
	HandCards     []poker.Card
	TableCards    []poker.Card
	ExpectedLevel poker.HandLevel
	ExpectedCards []poker.Card
}

var gameTests = []GameTest{
	GameTest{h("10S", "JS"), h("2H", "QS", "6D", "KS", "AS"), hl("StraightFlush", "A"), h("10S", "JS", "QS", "KS", "AS")},
	GameTest{h("AH", "3H"), h("2H", "4H", "5H", "2C", "3C"), hl("StraightFlush", "5"), h("AH", "2H", "3H", "4H", "5H")},
	GameTest{h("10S", "JS"), h("10C", "QD", "10D", "10H", "2S"), hl("FourOfAKind", "10", "Q"), h("10S", "10C", "10D", "10H", "QD")},
	GameTest{h("10S", "10C"), h("10D", "10H", "JD", "QD", "KD"), hl("FourOfAKind", "10", "K"), h("10S", "10C", "10D", "10H", "KD")},
	GameTest{h("5C", "10S"), h("5S", "5D", "5H", "AS", "KC"), hl("FourOfAKind", "5", "A"), h("5C", "5S", "5D", "5H", "AS")},
	GameTest{h("2S", "2H"), h("3D", "3H", "3C", "QD", "KS"), hl("FullHouse", "3", "2"), h("2S", "2H", "3H", "3C", "3D")},
	GameTest{h("2S", "3S"), h("4C", "4S", "2D", "2H", "3C"), hl("FullHouse", "2", "4"), h("2S", "2D", "2H", "4C", "4S")},
	GameTest{h("6H", "8H"), h("9H", "10H", "2H", "3S", "7C"), hl("Flush", "10", "9", "8", "6", "2"), h("10H", "9H", "8H", "6H", "2H")},
	GameTest{h("6S", "8H"), h("9H", "10H", "JH", "QH", "7H"), hl("StraightFlush", "Q"), h("QH", "JH", "10H", "9H", "8H")},
	GameTest{h("AS", "JH"), h("QC", "KD", "10S", "2C", "3C"), hl("Straight", "A"), h("AS", "KD", "QC", "JH", "10S")},
	GameTest{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), hl("Straight", "5"), h("AS", "2C", "3H", "4C", "5D")},
	GameTest{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), hl("ThreeOfAKind", "6", "K", "J"), h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), hl("ThreeOfAKind", "6", "K", "J"), h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), hl("TwoPair", "6", "4", "A"), h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), hl("TwoPair", "6", "4", "A"), h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("AS", "AH"), h("2S", "4C", "6D", "8S", "10D"), hl("OnePair", "A", "10", "8", "6"), h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("AS", "2S"), h("AH", "4C", "6D", "8S", "10D"), hl("OnePair", "A", "10", "8", "6"), h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("2S", "4S"), h("5D", "7S", "8S", "QH", "KH"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("2S", "KH"), h("5D", "7S", "8S", "QH", "4S"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("8S", "KH"), h("5D", "7S", "2S", "QH", "4S"), hl("HighCard", "K", "Q", "8", "7", "5"), h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("KS", "8C"), h("QH", "10C", "9H", "7H", "6S"), hl("Straight", "10"), h("10C", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "QC", "9H", "7H", "6S"), hl("Straight", "10"), h("10H", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), hl("Straight", "10"), h("10H", "9H", "8C", "7H", "6S")},
	// Uses none of the hand cards
	GameTest{h("2S", "3S"), h("AH", "10H", "JH", "QH", "KH"), hl("StraightFlush", "A"), h("AH", "KH", "QH", "JH", "10H")},
	// Uses one of the hand cards
	GameTest{h("AH", "3S"), h("2S", "10H", "JH", "QH", "KH"), hl("StraightFlush", "A"), h("AH", "KH", "QH", "JH", "10H")},
}

func TestHoldem(t *testing.T) {
	playerCount := 4
	randGen := rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests
	for i, testCase := range gameTests {
		if len(testCase.TableCards) != 5 {
			t.Fatalf("Wrong number of table cards %v", i)
		}
		pack := poker.NewPack()
		shuffleFixing(&pack, testCase.TableCards, testCase.HandCards, randGen)
		tableCards, playerCards := Deal(&pack, playerCount)
		outcomes := DealOutcomes(tableCards, playerCards)
		if len(tableCards) != 5 {
			t.Fatalf("Expected 5 table cards, found %v", len(tableCards))
		}
		for j, c := range tableCards {
			if c != testCase.TableCards[j] {
				t.Errorf("Expected %v at tableCards[%v], found %v", testCase.TableCards[j], j, c)
			}
		}
		if len(playerCards) != playerCount {
			t.Fatalf("Expected %v sets of player cards, found %v", playerCount, len(playerCards))
		}
		for j, c := range playerCards[0] {
			if c != testCase.HandCards[j] {
				t.Errorf("Expected %v at playerCards[0][%v], found %v", testCase.HandCards[j], j, c)
			}
		}
		if len(outcomes) != playerCount {
			t.Fatalf("Expected %v outcomes, found %v", playerCount, len(outcomes))
		}
		foundPlayer := false
		for _, outcome := range outcomes {
			if outcome.Player != 1 {
				continue
			}
			foundPlayer = true
			if !levelsEqual(outcome.Level, testCase.ExpectedLevel) {
				t.Errorf("Expected %v for case %v, found %v", testCase.ExpectedLevel, i, outcome.Level)
			}
			if len(outcome.Cards) != 5 {
				t.Fatalf("Expected 5 cards for case %v, found %v", i, len(outcome.Cards))
			}
			if !cardsEqual(testCase.ExpectedCards, outcome.Cards) {
				t.Errorf("Expected cards %q for case %v, found %q", testCase.ExpectedCards, i, outcome.Cards)
			}
		}
		if !foundPlayer {
			t.Errorf("Could not find player 1 in outcomes")
		}
		for j := 1; j < len(outcomes); j++ {
			for k := 0; k < j; k++ {
				prevLevel := outcomes[k].Level
				thisLevel := outcomes[j].Level
				if poker.Beats(thisLevel, prevLevel) {
					t.Errorf("%v beats %v but came after it", thisLevel, prevLevel)
				}
				if !poker.Beats(thisLevel, prevLevel) && !poker.Beats(prevLevel, thisLevel) && outcomes[k].Player > outcomes[j].Player {
					t.Errorf("%v and %v have equal hands (%v, %v) but %v came first", outcomes[k].Player, outcomes[j].Player, outcomes[k].Cards, outcomes[j].Cards, outcomes[k].Player)
				}
			}
		}
	}
}
