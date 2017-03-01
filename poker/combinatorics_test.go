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
	"sort"
	"testing"
)

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

type binomTestCase struct {
	n, r, expected int
}

func TestBinomial(t *testing.T) {
	cases := []binomTestCase{
		binomTestCase{3, 0, 1},
		binomTestCase{3, 1, 3},
		binomTestCase{3, 2, 3},
		binomTestCase{3, 3, 1},
		binomTestCase{4, 0, 1},
		binomTestCase{4, 1, 4},
		binomTestCase{4, 2, 6},
		binomTestCase{4, 3, 4},
		binomTestCase{4, 4, 1},
		binomTestCase{45, 2, 990},
	}

	for _, test := range cases {
		actual := binomial(test.n, test.r)
		if actual != test.expected {
			t.Errorf("Expected binomial(%v, %v) = %v, found %v", test.n, test.r, test.expected, actual)
		}
	}
}

func TestCardCombinations(t *testing.T) {
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

	choices := AllCardCombinations(cards, 3)
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
