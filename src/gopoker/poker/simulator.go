/*
Copyright 2013 Andrew Medworth

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

type Simulator struct {
	HandCount           int
	WinCount            int
	OurClassCounts      []int
	OpponentClassCounts []int
}

func (s *Simulator) SimulateHoldem(yourCards, tableCards []Card, players, handsToPlay int) {
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.OurClassCounts = make([]int, MAX_HANDCLASS)
	s.OpponentClassCounts = make([]int, MAX_HANDCLASS)

	p := NewPack()
	for i := 0; i < handsToPlay; i++ {
		p.shuffleFixing(tableCards, yourCards)
		won, ourLevel, opponentLevel := p.SimulateOneHoldemHand(players)
		if won {
			s.WinCount++
		}
		s.OurClassCounts[ourLevel.Class]++
		s.OpponentClassCounts[opponentLevel.Class]++
	}
}
