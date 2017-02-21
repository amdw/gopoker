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

// Compute all unique subsets of a set of cards, of a given size.
func allCardCombinations(pack []Card, numRequired int) [][]Card {
	// N choose 0 = 1 regardless of N/K
	if numRequired == 0 {
		return [][]Card{[]Card{}}
	}
	// N choose N = 1 regardless of N
	if numRequired == len(pack) {
		return [][]Card{pack}
	}
	// N choose K = N-1 choose K + N-1 choose K-1
	withoutFirst := allCardCombinations(pack[1:], numRequired)
	smallerWithoutFirst := allCardCombinations(pack[1:], numRequired-1)
	result := make([][]Card, len(withoutFirst)+len(smallerWithoutFirst))
	for i, sub := range smallerWithoutFirst {
		subset := make([]Card, len(sub)+1)
		subset[0] = pack[0]
		copy(subset[1:], sub)
		result[i] = subset
	}
	copy(result[len(smallerWithoutFirst):], withoutFirst)
	return result
}
