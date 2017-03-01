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
)

type Pack struct {
	Cards [52]Card
}

func (p *Pack) initialise() {
	i := 0
	for s := 0; s < 4; s++ {
		for r := 0; r < 13; r++ {
			p.Cards[i] = Card{Rank(r), Suit(s)}
			i++
		}
	}
}

// Shuffle the pack
func (p *Pack) Shuffle(randGen *rand.Rand) {
	for i := 0; i < 52; i++ {
		j := randGen.Intn(52-i) + i
		p.Cards[i], p.Cards[j] = p.Cards[j], p.Cards[i]
	}
}

func NewPack() Pack {
	var result Pack
	result.initialise()
	return result
}
