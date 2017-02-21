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

// A sub-function with an extra argument startSkippingAt to avoid duplication of results
func allChoicesSkipping(cards []Card, num, startSkippingAt int) [][]Card {
	if num >= len(cards) {
		return [][]Card{cards}
	}

	result := [][]Card{}

	// Call the function recursively with every possible one-smaller combination, starting skipping at the appropriate location.
	for i := startSkippingAt; i < len(cards); i++ {
		nextSmaller := make([]Card, len(cards)-1)
		j := 0
		for k, c := range cards {
			if k == i {
				continue
			}
			nextSmaller[j] = c
			j++
		}
		subChoices := allChoicesSkipping(nextSmaller, num, i)
		if len(subChoices) > 0 {
			newResult := make([][]Card, len(result)+len(subChoices))
			copy(newResult[0:len(result)], result)
			copy(newResult[len(result):], subChoices)
			result = newResult
		}
	}

	return result
}

// All ways to choose n cards from a set
func allChoices(cards []Card, num int) [][]Card {
	return allChoicesSkipping(cards, num, 0)
}
