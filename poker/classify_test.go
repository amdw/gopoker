/*
Copyright 2013, 2015-2017 Andrew Medworth

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
	"reflect"
	"testing"
)

type levelTest struct {
	l1        HandLevel
	l2        HandLevel
	isGreater bool
	isLess    bool
}

var levelTests = []levelTest{
	{hl("StraightFlush", "A"), hl("StraightFlush", "A"), false, false},
	{hl("StraightFlush", "A"), hl("StraightFlush", "K"), true, false},
	{hl("FourOfAKind", "9", "10"), hl("StraightFlush", "2"), false, true},
	{MinLevel(), hl("HighCard", "2", "2", "2", "2", "3"), false, true},
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

// We don't test regular classification directly here as it is done indirectly by game-specific tests

type classificationTest struct {
	cards         []Card
	expectedLevel HandLevel
}

var aceToFiveLowClassTests = []classificationTest{
	{h("8S", "5S", "4S", "3S", "2S"), hl("HighCard", "8", "5", "4", "3", "2")},
	{h("8C", "5S", "4S", "3S", "2S"), hl("HighCard", "8", "5", "4", "3", "2")},
	{h("5D", "4H", "3H", "2H", "AH"), hl("HighCard", "5", "4", "3", "2", "A")},
	{h("AS", "AH", "9S", "5S", "3S"), hl("OnePair", "A", "9", "5", "3")},
	{h("9S", "AS", "9C", "5S", "3S"), hl("OnePair", "9", "5", "3", "A")},
	{h("AC", "2C", "AH", "2H", "9D"), hl("TwoPair", "2", "A", "9")},
	{h("5D", "4H", "5S", "5C", "KD"), hl("ThreeOfAKind", "5", "K", "4")},
	{h("5D", "AH", "5S", "5C", "KD"), hl("ThreeOfAKind", "5", "K", "A")},
	{h("5D", "KS", "KD", "5C", "5H"), hl("FullHouse", "5", "K")},
	{h("AD", "KS", "KD", "AC", "AH"), hl("FullHouse", "A", "K")},
	{h("JS", "2D", "JC", "JH", "JD"), hl("FourOfAKind", "J", "2")},
	{h("AS", "2D", "AC", "AH", "AD"), hl("FourOfAKind", "A", "2")},
	{h("5S", "5C", "5H", "5D", "AS"), hl("FourOfAKind", "5", "A")},
}

func TestAceToFiveLowClassification(t *testing.T) {
	for _, test := range aceToFiveLowClassTests {
		level := ClassifyAceToFiveLow(test.cards)
		if !reflect.DeepEqual(level, test.expectedLevel) {
			t.Errorf("Expected %q as ace-to-five low classification for %q, found %q", test.expectedLevel, test.cards, level)
		}
	}
}

var aceToFiveLowBeatsTests = []levelTest{
	{hl("HighCard", "5", "4", "3", "2", "A"), hl("HighCard", "6", "5", "4", "3", "A"), true, false},
	{hl("HighCard", "6", "5", "4", "3", "2"), hl("HighCard", "5", "4", "3", "2", "A"), false, true},
	{hl("OnePair", "A", "4", "3", "2"), hl("HighCard", "8", "7", "6", "5", "4"), false, true},
	{hl("OnePair", "K", "6", "5", "4"), hl("OnePair", "A", "6", "5", "4"), false, true},
	{hl("OnePair", "K", "3", "2", "A"), hl("OnePair", "K", "4", "3", "2"), true, false},
	{hl("OnePair", "K", "4", "3", "A"), hl("OnePair", "K", "5", "4", "A"), true, false},
	{hl("TwoPair", "K", "J", "8"), hl("OnePair", "K", "J", "8", "7"), false, true},
	{hl("TwoPair", "K", "J", "8"), hl("TwoPair", "K", "A", "8"), false, true},
	{hl("TwoPair", "K", "Q", "9"), hl("ThreeOfAKind", "3", "6", "5"), true, false},
	{hl("ThreeOfAKind", "3", "6", "5"), hl("ThreeOfAKind", "2", "6", "5"), false, true},
	{hl("ThreeOfAKind", "3", "6", "5"), hl("ThreeOfAKind", "3", "6", "5"), false, false},
	{hl("FullHouse", "A", "2"), hl("ThreeOfAKind", "3", "6", "5"), false, true},
	{hl("FullHouse", "A", "2"), hl("FullHouse", "2", "A"), true, false},
	{hl("FourOfAKind", "A", "2"), hl("FullHouse", "A", "2"), false, true},
	{hl("FourOfAKind", "A", "2"), hl("FourOfAKind", "2", "A"), true, false},
}

func TestBeatsAceToFiveLow(t *testing.T) {
	for _, test := range aceToFiveLowBeatsTests {
		gt := BeatsAceToFiveLow(test.l1, test.l2)
		lt := BeatsAceToFiveLow(test.l2, test.l1)
		if gt != test.isGreater {
			t.Errorf("Expected %q > %q = %v at ace-to-five low, found %v", test.l1, test.l2, test.isGreater, gt)
		}
		if lt != test.isLess {
			t.Errorf("Expected %q < %q = %v at ace-to-five low, found %v", test.l1, test.l2, test.isLess, lt)
		}
	}
}
