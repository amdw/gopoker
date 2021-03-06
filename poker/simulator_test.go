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
	"math"
	"testing"
)

func TestPotOdds(t *testing.T) {
	sim := Simulator{}
	sim.HandCount = 10000
	tests := map[float64]float64{0.0: 0.0, float64(sim.HandCount) / 2.0: 1.0, float64(sim.HandCount): math.Inf(1)}
	for potsWon, expected := range tests {
		sim.PotsWon = potsWon
		breakEven := sim.PotOddsBreakEven()
		if breakEven != expected {
			t.Errorf("Expected pot odds break even %v, found %v", expected, breakEven)
		}
	}
}
