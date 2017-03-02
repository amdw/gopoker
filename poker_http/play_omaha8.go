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
	"fmt"
	"github.com/amdw/gopoker/omaha8"
	"github.com/amdw/gopoker/poker"
	"math/rand"
	"net/http"
	"time"
)

func printTickCell(w http.ResponseWriter, tick bool) {
	if tick {
		fmt.Fprintf(w, `<td class="tickcell">&#9989;</td>`)
	} else {
		fmt.Fprintf(w, "<td></td>")
	}
}

func PlayOmaha8(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	players, err := getPlayers(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting player count: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, `<html lang="en">`)
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, `<meta charset="utf-8">`)
	fmt.Fprintln(w, `<meta http-equiv="X-UA-Compatible" content="IE=edge">`)
	fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(w, `<title>Example game of Omaha/8</title>`)
	fmt.Fprintln(w, `<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">`)
	fmt.Fprintln(w, "<style>")
	fmt.Fprintln(w, "th { text-align: center }")
	fmt.Fprintln(w, "td.numcell { text-align: right }")
	fmt.Fprintln(w, ".nothing { color: lightgray }")
	fmt.Fprintln(w, "td.tickcell { text-align: center }")
	fmt.Fprintln(w, "</style>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, `<div class="container-fluid">`)

	fmt.Fprintln(w, "<h1>Example Omaha/8 game</h1>")

	fmt.Fprintln(w, `<form method="get">`)
	fmt.Fprintln(w, `<div class="form-group"><label for="playerCount">Players</label>`)
	fmt.Fprintf(w, `<input type="text" id="playerCount" name="%v" value="%v" class="form-control"/>`, playersKey, players)
	fmt.Fprintln(w, "</div>")
	fmt.Fprintln(w, `<button type="submit" class="btn btn-default">Rerun</button></form>`)

	pack := poker.NewPack()
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	pack.Shuffle(randGen)
	tableCards, playerCards := omaha8.Deal(&pack, players)
	playerOutcomes := omaha8.PlayerOutcomes(tableCards, playerCards)

	fmt.Fprintf(w, "<h3>Table cards</h3><p>%v</p>", formatCards(tableCards))
	fmt.Fprintln(w, "<h3>Player cards</h3><ul>")
	for playerIdx := 0; playerIdx < len(playerCards); playerIdx++ {
		fmt.Fprintf(w, "<li>Player %v: %v</li>\n", playerIdx+1, formatCards(playerCards[playerIdx]))
	}
	fmt.Fprintln(w, "</ul>")

	fmt.Fprintln(w, "<h3>Results</h3>")
	fmt.Fprintln(w, `<table class="table table-bordered">`)
	fmt.Fprintln(w, `<tr><th rowspan="2">Player</th><th colspan="3">High</th><th colspan="3">Low</th><th rowspan="2">Winnings</th></tr>`)
	fmt.Fprintln(w, `<tr><th>Hand</th><th>Cards</th><th>Win?</th><th>Hand</th><th>Cards</th><th>Win?</th></tr>`)

	for _, outcome := range playerOutcomes {
		fmt.Fprintln(w, "<tr>")
		fmt.Fprintf(w, `<td>%v</td>`, outcome.Player)
		fmt.Fprintf(w, `<td>%v</td><td>%v</td>`, outcome.Level.HighLevel.PrettyPrint(), formatCards(outcome.Level.HighHand))
		printTickCell(w, outcome.IsHighWinner)
		if outcome.Level.LowLevelQualifies {
			fmt.Fprintf(w, `<td>%v</td><td>%v</td>`, outcome.Level.LowLevel.PrettyPrint(), formatCards(outcome.Level.LowHand))
		} else {
			fmt.Fprintf(w, `<td class="nothing">None</td><td class="nothing">None</td>`)
		}
		printTickCell(w, outcome.IsLowWinner)
		fracClass := ""
		if outcome.PotFractionWon() == 0 {
			fracClass = " nothing"
		}
		fmt.Fprintf(w, `<td class="numcell%v">%.1f%%</td>`, fracClass, 100*outcome.PotFractionWon())
		fmt.Fprintln(w, "</tr>")
	}

	fmt.Fprintln(w, "</table>")

	fmt.Fprintln(w, "</div>")
	fmt.Fprintln(w, "</body></html>")
}
