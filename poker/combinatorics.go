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

// Compute each unique subset of a set of cards of a given size, and send it down a channel.
// This implementation is experimental; one might expect it to perform better than the recursive
// implementation because of the lack of repeated work and reduced memory allocation, but in fact
// what we gain in these areas seems to be almost exactly balanced by the cost of semaphore operations.
func enumerateCardCombinations(pack []Card, numRequired int, c chan []Card) {
	indices := make([]int, numRequired)
	// Standard algorithm to enumerate k-combinations
	for i := 0; i < numRequired; i++ {
		indices[i] = i
	}
	for {
		// Construct and send a combination
		combination := make([]Card, numRequired)
		for i := 0; i < numRequired; i++ {
			combination[i] = pack[indices[i]]
		}
		c <- combination

		// Advance to the next combination
		if indices[numRequired-1] < len(pack)-1 {
			indices[numRequired-1]++
		} else {
			i := numRequired - 1
			for i >= 0 && indices[i] == i+len(pack)-numRequired {
				i--
			}
			if i < 0 {
				close(c)
				break
			}
			indices[i]++
			for j := i + 1; j < numRequired; j++ {
				indices[j] = indices[j-1] + 1
			}
		}
	}
}
