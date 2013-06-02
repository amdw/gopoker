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
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func TestCardBasics(t *testing.T) {
	for s := Heart; s <= Club; s++ {
		for r := Two; r <= Ace; r++ {
			cs := fmt.Sprintf("%v%v", r, s)
			c1 := C(cs)
			c2 := Card{r, s}
			if c1 != c2 {
				t.Errorf("Expected %q, found %q", c2, c1)
			}
			if c1.String() != cs {
				t.Errorf("Expected %q, found %q", cs, c1.String())
			}
		}
	}
	// Test lower-case conversion
	if C("JS") != C("js") {
		t.Errorf("Should be able to accept lower-case cards, but found %v vs %v", C("JS"), C("js"))
	}
}

type LevelTest struct {
	l1        HandLevel
	l2        HandLevel
	isGreater bool
	isLess    bool
}

func hl(class HandClass, tiebreaks []Rank) HandLevel {
	return HandLevel{class, tiebreaks}
}

var levelTests = []LevelTest{
	{hl(StraightFlush, []Rank{Ace}), hl(StraightFlush, []Rank{Ace}), false, false},
	{hl(StraightFlush, []Rank{Ace}), hl(StraightFlush, []Rank{King}), true, false},
	{hl(FourOfAKind, []Rank{Nine, Ten}), hl(StraightFlush, []Rank{Two}), false, true},
}

func TestLevels(t *testing.T) {
	for _, ltst := range levelTests {
		gt := Beats(ltst.l1, ltst.l2)
		lt := Beats(ltst.l2, ltst.l1)
		if gt != ltst.isGreater {
			t.Errorf("Expected %q beats %q == %v, found %v", ltst.l1, ltst.l2, ltst.isGreater, gt)
		}
		if lt != ltst.isLess {
			t.Errorf("Expected %q beats %q == %v, found %v", ltst.l2, ltst.l1, ltst.isLess, lt)
		}
	}
}

type ClassTest struct {
	mandatory     []Card
	optional      []Card
	expectedLevel HandLevel
	expectedCards []Card
}

func h(cards ...string) []Card {
	result := make([]Card, len(cards))
	for i, c := range cards {
		result[i] = C(c)
	}
	return result
}

var classTests = []ClassTest{
	ClassTest{h("AS", "KS", "QS", "JS", "10S"), h(), HandLevel{StraightFlush, []Rank{Ace}}, h("AS", "KS", "QS", "JS", "10S")},
	ClassTest{h("9D", "10S", "9S", "9H", "9C"), h(), HandLevel{FourOfAKind, []Rank{Nine, Ten}}, h("9D", "10S", "9S", "9H", "9C")},
	ClassTest{h("10S", "JS"), h("2H", "QS", "6D", "KS", "AS"), HandLevel{StraightFlush, []Rank{Ace}}, h("10S", "JS", "QS", "KS", "AS")},
	ClassTest{h("10S", "JS"), h("10C", "QD", "10D", "10H", "2S"), HandLevel{FourOfAKind, []Rank{Ten, Jack}}, h("10S", "10C", "10D", "10H", "JS")},
	ClassTest{h("10S", "10C"), h("10D", "10H", "JD", "QD", "KD", "AD"), HandLevel{FourOfAKind, []Rank{Ten, Ace}}, h("10S", "10C", "10D", "10H", "AD")},
	ClassTest{h("2S", "2H"), h("3D", "3H", "3C", "QD", "KS"), HandLevel{FullHouse, []Rank{Three, Two}}, h("2S", "2H", "3H", "3C", "3D")},
	ClassTest{h("2S", "3S"), h("4H", "4D", "4C", "4S", "2D", "2H", "3C"), HandLevel{FullHouse, []Rank{Two, Three}}, h("2S", "2D", "2H", "3S", "3C")},
	ClassTest{h("6H", "8H"), h("9H", "10H", "2H", "3S", "7C"), HandLevel{Flush, []Rank{Ten, Nine, Eight, Six, Two}}, h("10H", "9H", "8H", "6H", "2H")},
	ClassTest{h("6S", "8H"), h("9H", "10H", "JH", "QH", "7H"), HandLevel{Straight, []Rank{Ten, Nine, Eight, Seven, Six}}, h("10H", "9H", "8H", "7H", "6S")},
	ClassTest{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), HandLevel{Straight, []Rank{Five, Four, Three, Two, Ace}}, h("AS", "2C", "3H", "4C", "5D")},
	ClassTest{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}}, h("6S", "6D", "6C", "KH", "JC")},
	ClassTest{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), HandLevel{ThreeOfAKind, []Rank{Six, King, Two}}, h("6S", "6D", "6C", "KH", "2S")},
	ClassTest{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	ClassTest{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	ClassTest{h("KS", "8C"), h("10H", "10C", "9H", "7H", "6S"), HandLevel{OnePair, []Rank{Ten, King, Nine, Eight}}, h("KS", "8C", "10H", "10C", "9H")},
	ClassTest{h("KS", "10H"), h("8C", "10C", "9H", "7H", "6S"), HandLevel{OnePair, []Rank{Ten, King, Nine, Eight}}, h("KS", "8C", "10H", "10C", "9H")},
	ClassTest{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), HandLevel{HighCard, []Rank{King, Ten, Nine, Eight, Seven}}, h("KS", "10H", "9H", "8C", "7H")},
}

func levelsEqual(l1, l2 HandLevel) bool {
	if l1.Class != l2.Class {
		return false
	}
	if len(l1.Tiebreaks) != len(l2.Tiebreaks) {
		return false
	}
	for i := 0; i < len(l1.Tiebreaks); i++ {
		if l1.Tiebreaks[i] != l2.Tiebreaks[i] {
			return false
		}
	}
	return true
}

func cardsEqual(c1, c2 []Card) bool {
	if len(c1) != len(c2) {
		return false
	}
	sort.Sort(CardSorter{c1, false})
	sort.Sort(CardSorter{c2, false})
	for i := 0; i < len(c1); i++ {
		if c1[i] != c2[i] {
			return false
		}
	}
	return true
}

func TestClassification(t *testing.T) {
	for _, ct := range classTests {
		level, cards := Classify(ct.mandatory, ct.optional)
		if !levelsEqual(ct.expectedLevel, level) {
			t.Errorf("Expected %q, found %q for %q / %q", ct.expectedLevel, level, ct.mandatory, ct.optional)
		}
		if !cardsEqual(ct.expectedCards, cards) {
			t.Errorf("Expected cards %q, found %q for %q / %q", ct.expectedCards, cards, ct.mandatory, ct.optional)
		}
	}
}

func TestBuildStraights(t *testing.T) {
	suits := [][]Suit{{Spade}, {Club, Heart}, {Spade}, {Spade}, {Spade}}
	expectedOutput := [][]Suit{{Spade, Club, Spade, Spade, Spade}, {Spade, Heart, Spade, Spade, Spade}}
	actualOutput := buildStraights(suits)
	for i := range expectedOutput {
		for j := range expectedOutput[i] {
			if expectedOutput[i][j] != actualOutput[i][j] {
				t.Errorf("Expected %q at %v but found %q", expectedOutput[i], i, actualOutput[i])
				break
			}
		}
	}
}

func TestFixedShuffle(t *testing.T) {
	pack := NewPack()
	myCards := h("KS", "AC")
	tableCards := h("10D", "2C", "AS")
	pack.randGen = rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests

	for i := 0; i < 100; i++ {
		pack.Shuffle()
		pack.shuffleFixing(tableCards, myCards)
		for j, c := range h("10D", "2C", "AS") {
			if pack.Cards[j] != c {
				t.Fatalf("Expected %v at Cards[%v], found %v", c, j, pack.Cards[j])
			}
		}
		for j, c := range h("KS", "AC") {
			if pack.Cards[j+5] != c {
				t.Fatalf("Expected %v at Cards[%v], found %v", c, j+5, pack.Cards[j+5])
			}
		}
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
}

func TestSorting(t *testing.T) {
	cards := h("AS", "JC", "QD", "3C", "4S", "10C")

	sort.Sort(CardSorter{cards, false})
	for i, c := range h("AS", "QD", "JC", "10C", "4S", "3C") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-high list, found %v", c, i, cards[i])
		}
	}

	sort.Sort(CardSorter{cards, true})
	for i, c := range h("QD", "JC", "10C", "4S", "3C", "AS") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-low list, found %v", c, i, cards[i])
		}
	}
}
