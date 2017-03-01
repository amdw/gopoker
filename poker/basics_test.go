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
	"fmt"
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

type rankOrderTest struct {
	r1, r2      Rank
	lessAceHigh bool
	lessAceLow  bool
}

var rankOrderTests = []rankOrderTest{
	{Two, Three, true, true},
	{Three, Two, false, false},
	{Jack, King, true, true},
	{King, Jack, false, false},
	{Ace, Two, false, true},
	{Two, Ace, true, false},
	{Ace, King, false, true},
	{King, Ace, true, false},
	{Two, Two, false, false},
	{Three, Three, false, false},
}

func TestRankOrdering(t *testing.T) {
	for _, test := range rankOrderTests {
		lah := IsRankLess(test.r1, test.r2, false)
		lal := IsRankLess(test.r1, test.r2, true)
		if lah != test.lessAceHigh {
			t.Errorf("Expected %v < %v = %v ace high, found %v", test.r1, test.r2, test.lessAceHigh, lah)
		}
		if lal != test.lessAceLow {
			t.Errorf("Expected %v < %v = %v ace low, found %v", test.r1, test.r2, test.lessAceLow, lal)
		}
	}
}

func TestSorting(t *testing.T) {
	cards := h("AS", "JC", "QD", "3C", "4S", "10C")

	SortCards(cards, false)
	for i, c := range h("AS", "QD", "JC", "10C", "4S", "3C") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-high list, found %v", c, i, cards[i])
		}
	}

	SortCards(cards, true)
	for i, c := range h("QD", "JC", "10C", "4S", "3C", "AS") {
		if cards[i] != c {
			t.Errorf("Expected %v at position %v of ace-low list, found %v", c, i, cards[i])
		}
	}
}
