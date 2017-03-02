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
package poker_http

import (
	"encoding/json"
	"fmt"
	"github.com/amdw/gopoker/omaha8"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func SimulateOmaha8(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	params, err := getSimulationParams(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get simulation parameters: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, `<html lang="en"><head><title>Omaha/8 Simulator</title></head><body>`)
	fmt.Fprintln(w, "<h1>Omaha/8 Simulator</h1>")

	//if len(params.tableCards) > 0 || len(params.yourCards) > 0 || params.forceComputation {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	simulator := omaha8.SimulateOmaha8(params.tableCards, params.yourCards, params.players, params.handsToPlay, randGen)

	fmt.Fprintln(w, "<h2>Results</h2>")

	breakEven := simulator.PotOddsBreakEven()
	if math.IsInf(breakEven, 1) {
		fmt.Fprintln(w, "<p><b>Any</b> bet has positive expected value! :)</p>")
	} else {
		fmt.Fprintf(w, "<p>A bet up to %.1f%% of the pot has positive expected value.</p>", 100.0*breakEven)
	}

	fmt.Fprintln(w, "<code>")
	json.NewEncoder(w).Encode(simulator)
	fmt.Fprintln(w, "</code>")
	//}

	fmt.Fprintln(w, "</body></html>")
}
