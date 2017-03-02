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
package poker

import (
	"math/rand"
	"testing"
)

func TestShuffle(t *testing.T) {
	pack := NewPack()
	randGen := rand.New(rand.NewSource(1234)) // Deterministic for predictable tests
	pack.Shuffle(randGen)
	TestPackPermutation(&pack, t)
}

func TestIndexOf(t *testing.T) {
	pack := NewPack()
	randGen := rand.New(rand.NewSource(1234)) // Deterministic for predictable tests
	tests := 1000

	card := C("AS")
	randCheck := make(map[int]int)
	checkLimit := 10
	for i := 0; i < tests; i++ {
		pack.Shuffle(randGen)
		idx := pack.IndexOf(card)
		if pack.Cards[idx] != card {
			t.Errorf("Expected %v at %v but found %v", card, idx, pack.Cards[idx])
		}
		if len(randCheck) < checkLimit {
			randCheck[idx]++
		}
	}
	if len(randCheck) < checkLimit {
		t.Errorf("Suspicious lack of randomness - only indices found: %q", randCheck)
	}
}
