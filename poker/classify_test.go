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
	"testing"
)

type LevelTest struct {
	l1        HandLevel
	l2        HandLevel
	isGreater bool
	isLess    bool
}

var hl = TestMakeHandLevel

var levelTests = []LevelTest{
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

// We don't test classification directly here as it is done indirectly by game-specific tests
