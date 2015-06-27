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
	"gopoker/poker"
	"log"
	"net/http"
	"strconv"
)

const rank1Key = "rank1"
const rank2Key = "rank2"
const sameSuitKey = "samesuit"
const handsToPlayKey = "handstoplay"

func StartingCards(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, `<!DOCTYPE html>`)
	fmt.Fprintln(w, `<html>`)
	fmt.Fprintln(w, `<head><title>Texas Hold'em starting cards</title>`)
	fmt.Fprintln(w, `<style>`)
	fmt.Fprintln(w, `table.resultTable { border-collapse: collapse }`)
	fmt.Fprintln(w, `th.resultTable, td.resultTable { border: 1px solid black; padding: 3px }`)
	fmt.Fprintln(w, `</style>`)
	fmt.Fprintln(w, `<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.3.16/angular.min.js"></script>`)
	fmt.Fprintln(w, `</head><body>`)

	fmt.Fprintln(w, `<div ng-app="startingCardsApp" ng-controller="StartingCardsController">`)
	fmt.Fprintln(w, `<h1>Texas Hold'em starting cards</h1>`)
	fmt.Fprintln(w, `<table>`)
	fmt.Fprintln(w, `<tr><td>Players</td><td><input type="text" ng-model="players"/></td></tr>`)
	fmt.Fprintln(w, `<tr><td>Hands to simulate</td><td><input type="text" ng-model="handsToPlay"/></td></tr>`)
	fmt.Fprintln(w, `<tr><td><button ng-click="simulate()" ng-disabled="started">Simulate</button></td><td>Simulation {{status()}}</td></tr>`)
	fmt.Fprintln(w, `</table>`)
	fmt.Fprintln(w, `<h2>Results</h2>`)
	fmt.Fprintln(w, `<table class="resultTable">`)
	fmt.Fprintln(w, `<tr><th rowspan="2" class="resultTable">Cards</th><th colspan="3" class="resultTable">Win probability (versus prior)</th></tr>`)
	fmt.Fprintln(w, `<tr><th class="resultTable">You</th><th class="resultTable">Best opponent</th><th class="resultTable">Random opponent</th></tr>`)
	fmt.Fprintln(w, `<tr ng-repeat="result in results">`)
	fmt.Fprintln(w, `<td class="resultTable">{{result.Cards}}</td>`)
	fmt.Fprintln(w, `<td class="resultTable">{{(100 * result.WinCount / result.HandCount) | number : 1}}% ({{result.WinCount / (result.HandCount / players) | number : 2}})</td>`)
	fmt.Fprintln(w, `<td class="resultTable">{{(100 * result.BestOpponentWinCount / result.HandCount) | number : 1}}% ({{result.BestOpponentWinCount / (result.HandCount * (players - 1) / players) | number : 2}})</td>`)
	fmt.Fprintln(w, `<td class="resultTable">{{(100 * result.RandomOpponentWinCount / result.HandCount) | number : 1}}% ({{result.RandomOpponentWinCount / (result.HandCount / players) | number : 2}})</td>`)
	fmt.Fprintln(w, `</tr>`)
	fmt.Fprintln(w, `</table>`)
	fmt.Fprintln(w, `</div>`)

	fmt.Fprintln(w, `<script>`)
	fmt.Fprintln(w, `var app = angular.module('startingCardsApp', []);`)

	fmt.Fprintln(w, `app.controller('StartingCardsController', function($scope, $http) {`)
	fmt.Fprintln(w, `$scope.players = 7;`)
	fmt.Fprintln(w, `$scope.handsToPlay = 10000;`)
	fmt.Fprintln(w, `$scope.results = [];`)
	fmt.Fprintln(w, `$scope.started = false;`)
	fmt.Fprintln(w, `$scope.requestsMade = 0;`)
	fmt.Fprintln(w, `$scope.resultsPending = 0;`)
	fmt.Fprintln(w, `$scope.errors = 0;`)
	fmt.Fprintln(w, `$scope.status = function() {`)
	fmt.Fprintln(w, `if (!$scope.started) { return "not started"; }`)
	fmt.Fprintln(w, `if ($scope.resultsPending > 0) {`)
	fmt.Fprintln(w, `return "in progress (" + $scope.resultsPending + " of " + $scope.requestsMade + " requests pending" + ($scope.errors > 0 ? ("; " + $scope.errors + " errors - see console for details") : "") + ")"; }`)
	fmt.Fprintln(w, `return "complete (reload page to restart)";`)
	fmt.Fprintln(w, `};`)
	fmt.Fprintln(w, `$scope.simulateOne = function(rank1, rank2, samesuit) {`)
	fmt.Fprintln(w, `var url = "/startingcards/sim?rank1=" + rank1 + "&rank2=" + rank2 + "&samesuit=" + samesuit + "&players=" + $scope.players + "&handstoplay=" + $scope.handsToPlay;`)
	fmt.Fprintln(w, `$http.get(url).success(function (response) {$scope.onResult(rank1, rank2, samesuit, response)}).error(function (response, status) {$scope.onError(rank1, rank2, samesuit, response, status)});`)
	fmt.Fprintln(w, `$scope.requestsMade += 1;`)
	fmt.Fprintln(w, `$scope.resultsPending += 1;`)
	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `$scope.simulate = function() {`)
	fmt.Fprintln(w, `$scope.started = true;`)
	fmt.Fprintln(w, `var ranks = ['A','2','3','4','5','6','7','8','9','10','J','Q','K'];`)
	fmt.Fprintln(w, `for (i in ranks) {`)
	fmt.Fprintln(w, `for (j in ranks) {`)
	fmt.Fprintln(w, `if (j > i) { break; } // Avoid duplicates`)
	fmt.Fprintln(w, `$scope.simulateOne(ranks[i], ranks[j], false);`)
	fmt.Fprintln(w, `if (i != j) { $scope.simulateOne(ranks[i], ranks[j], true); }`)
	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `};`)
	fmt.Fprintln(w, `$scope.onResult = function (rank1, rank2, samesuit, result) {`)
	fmt.Fprintln(w, `result.Cards = rank1 + rank2 + (samesuit ? "s" : "");`)
	fmt.Fprintln(w, `$scope.results.push(result);`)
	fmt.Fprintln(w, `$scope.results.sort(function(a,b) {return b.WinCount - a.WinCount});`)
	fmt.Fprintln(w, `$scope.resultsPending -= 1;`)
	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `$scope.onError = function (rank1, rank2, samesuit, error, status) {`)
	fmt.Fprintln(w, `$scope.errors += 1;`)
	fmt.Fprintln(w, `console.log("Got error for " + rank1 + rank2 + samesuit + ": " + status + "\n" + JSON.stringify(error));`)
	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `});`)
	fmt.Fprintln(w, `</script>`)

	fmt.Fprintln(w, `</body></html>`)
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

func (params SimParams) RunSimulation() poker.Simulator {
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
