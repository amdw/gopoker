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
package poker_http

import (
	"encoding/json"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

const rank1Key = "rank1"
const rank2Key = "rank2"
const sameSuitKey = "samesuit"
const handsToPlayKey = "handstoplay"

func StartingCards(staticBaseDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		path := path.Join(staticBaseDir, "starting_cards.html")
		file, err := os.Open(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not load %v: %v", path, err), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not write %v: %v", path, err), http.StatusInternalServerError)
			return
		}
	}
}

func getRank(req *http.Request, key string, w http.ResponseWriter) (poker.Rank, bool) {
	var rank poker.Rank
	var err error
	if rstrs, ok := req.Form[key]; ok && len(rstrs) > 0 {
		rank, err = poker.MakeRank(rstrs[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad rank '%v': %v", rstrs[0], err), http.StatusBadRequest)
			return 0, false
		}
	} else {
		http.Error(w, fmt.Sprintf("Missing required %v", key), http.StatusBadRequest)
		return 0, false
	}
	return rank, true
}

type SimParams struct {
	StartingPair poker.StartingPair
	Players      int
	HandsToPlay  int
}

func (params SimParams) RunSimulation() *poker.Simulator {
	return params.StartingPair.RunSimulation(params.Players, params.HandsToPlay)
}

func getStartingPair(req *http.Request, w http.ResponseWriter) (SimParams, bool) {
	req.ParseForm()
	rank1, ok := getRank(req, rank1Key, w)
	if !ok {
		return SimParams{}, false
	}
	rank2, ok := getRank(req, rank2Key, w)
	if !ok {
		return SimParams{}, false
	}
	sameSuit := false
	var err error
	if ssstrs, ok := req.Form[sameSuitKey]; ok && len(ssstrs) > 0 {
		sameSuit, err = strconv.ParseBool(ssstrs[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad %v: %v", sameSuitKey, err), http.StatusBadRequest)
			return SimParams{}, false
		}
	}
	players := 7
	if pstrs, ok := req.Form[playersKey]; ok && len(pstrs) > 0 {
		var players64 int64
		players64, err = strconv.ParseInt(pstrs[0], 10, 31)
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad %v: %v", playersKey, err), http.StatusBadRequest)
			return SimParams{}, false
		}
		players = int(players64)
	}
	handsToPlay := 10000
	if htpstrs, ok := req.Form[handsToPlayKey]; ok && len(htpstrs) > 0 {
		var handsToPlay64 int64
		handsToPlay64, err = strconv.ParseInt(htpstrs[0], 10, 31)
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad %v: %v", handsToPlayKey, err), http.StatusBadRequest)
			return SimParams{}, false
		}
		handsToPlay = int(handsToPlay64)
	}
	startingPair := poker.StartingPair{rank1, rank2, sameSuit}
	err = startingPair.Validate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad pair: %v", err), http.StatusBadRequest)
		return SimParams{}, false
	}
	return SimParams{startingPair, players, handsToPlay}, true
}

func SimulateStartingCards(w http.ResponseWriter, req *http.Request) {
	simParams, ok := getStartingPair(req, w)
	if !ok {
		return
	}
	log.Println("Simulating", simParams)
	simulator := simParams.RunSimulation()
	log.Println("Simulation", simParams, "complete")
	json.NewEncoder(w).Encode(simulator)
}
