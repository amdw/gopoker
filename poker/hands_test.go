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

func hl(class HandClass, tieBreaks ...Rank) HandLevel {
	return HandLevel{class, tieBreaks}
}

var levelTests = []LevelTest{
	{hl(StraightFlush, Ace), hl(StraightFlush, Ace), false, false},
	{hl(StraightFlush, Ace), hl(StraightFlush, King), true, false},
	{hl(FourOfAKind, Nine, Ten), hl(StraightFlush, Two), false, true},
	{MinLevel(), hl(HighCard, Two, Two, Two, Two, Three), false, true},
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

func h(cards ...string) []Card {
	result := make([]Card, len(cards))
	for i, c := range cards {
		result[i] = C(c)
	}
	return result
}

func TestSorting(t *testing.T) {
	cards := h("AS", "JC", "QD", "3C", "4S", "10C")

	sortCards(cards, false)
	for i, c := range h("AS", "QD", "JC", "10C", "4S", "3C") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-high list, found %v", c, i, cards[i])
		}
	}

	sortCards(cards, true)
	for i, c := range h("QD", "JC", "10C", "4S", "3C", "AS") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-low list, found %v", c, i, cards[i])
		}
	}
}

func lexSortHands(hands [][]Card) {
	sort.Slice(hands, func(i, j int) bool {
		for k := 0; k < len(hands[i]) && k < len(hands[j]); k++ {
			c1, c2 := hands[i][k], hands[j][k]
			if c1 != c2 {
				if c1.Rank != c2.Rank {
					return c1.Rank < c2.Rank
				}
				return c1.Suit < c2.Suit
			}
		}
		return false
	})
}

func TestAllChoices(t *testing.T) {
	cards := h("AS", "QD", "JC", "3C", "2H")

	expectedChoices := [][]Card{
		h("AS", "QD", "JC"),
		h("AS", "QD", "3C"),
		h("AS", "JC", "3C"),
		h("QD", "JC", "3C"),
		h("AS", "QD", "2H"),
		h("AS", "JC", "2H"),
		h("QD", "JC", "2H"),
		h("AS", "3C", "2H"),
		h("QD", "3C", "2H"),
		h("JC", "3C", "2H"),
	}
	lexSortHands(expectedChoices)

	choices := allChoices(cards, 3)
	lexSortHands(choices)

	if len(choices) != len(expectedChoices) {
		t.Errorf("Expected %v choices, found %v: %v / %v", len(expectedChoices), len(choices), expectedChoices, choices)
	}

	for i := range choices {
		if len(expectedChoices[i]) != len(choices[i]) {
			t.Errorf("Expected %v at choices[%v], found %v", expectedChoices[i], i, choices[i])
		}
		for j := range expectedChoices[i] {
			if expectedChoices[i][j] != choices[i][j] {
				t.Errorf("Expected %v at choices[%v], found %v", expectedChoices[i], i, choices[i])
				break
			}
		}
	}
}
