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

func TestMakeHand(cards ...string) []Card {
	result := make([]Card, len(cards))
	for i, c := range cards {
		result[i] = C(c)
	}
	return result
}

func parseHandClass(handClassStr string) HandClass {
	switch handClassStr {
	case "StraightFlush":
		return StraightFlush
	case "FourOfAKind":
		return FourOfAKind
	case "FullHouse":
		return FullHouse
	case "Flush":
		return Flush
	case "Straight":
		return Straight
	case "ThreeOfAKind":
		return ThreeOfAKind
	case "TwoPair":
		return TwoPair
	case "OnePair":
		return OnePair
	case "HighCard":
		return HighCard
	default:
		panic(fmt.Sprintf("Unknown hand class %v", handClassStr))
	}
}

func TestMakeHandLevel(handClassStr string, tieBreakRankStrs ...string) HandLevel {
	class := parseHandClass(handClassStr)
	tieBreaks := make([]Rank, len(tieBreakRankStrs))
	for i, rankStr := range tieBreakRankStrs {
		var err error
		tieBreaks[i], err = MakeRank(rankStr)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse rank %v", rankStr))
		}
	}
	return HandLevel{class, tieBreaks}
}

// Assert that the pack contains exactly one of every card
func TestPackPermutation(pack *Pack, t *testing.T) {
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
