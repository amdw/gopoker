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
package poker_http

import (
	"fmt"
	"github.com/amdw/gopoker/holdem"
	"github.com/amdw/gopoker/poker"
	"math/rand"
	"net/http"
	"time"
)

func PlayHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	players, err := getPlayers(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting player count: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "<html><head><title>A game of Texas Hold'em</title>")
	fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(w, "</head><body><h1>A game of Texas Hold'em</h1>")

	fmt.Fprintf(w, `<form method="get">Players: <input type="text" name="%v" value="%v"/><input type="submit" value="Rerun"/></form>`, playersKey, players)

	pack := poker.NewPack()
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	pack.Shuffle(randGen)
	onTable, playerCards := holdem.Deal(&pack, players)
	sortedOutcomes := holdem.DealOutcomes(onTable, playerCards)
	fmt.Fprintf(w, "<h2>Table cards</h2><p>%v</p>", formatCards(onTable))
	fmt.Fprintf(w, "<h2>Player cards</h2><ul>")
	for player := 0; player < players; player++ {
		fmt.Fprintf(w, "<li>Player %v: %v</li>", player+1, formatCards(playerCards[player]))
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "<h2>Results</h2><table><tr><th>Position</th><th>Player</th><th>Hand</th><th>Cards</th></tr>")
	for i, outcome := range sortedOutcomes {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", i+1, outcome.Player, outcome.Level.PrettyPrint(), formatCards(outcome.Cards))
	}
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</body></html>")
}
