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
	"errors"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Menu(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "<html><head>")
	fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(w, "<title>Poker</title></head><body><h1>Poker</h1><ul>")
	fmt.Fprintln(w, `<li><a href="/play">Play</a></li>`)
	fmt.Fprintln(w, `<li><a href="/simulate">Simulate</a></li>`)
	fmt.Fprintln(w, `<li><a href="/startingcards">Starting cards</a></li>`)
	fmt.Fprintln(w, "</ul></body></html>")
}

func formatCards(cards []poker.Card) string {
	cardStrings := make([]string, len(cards))
	for i, c := range cards {
		cardStrings[i] = c.HTML()
	}
	return strings.Join(cardStrings, ", ")
}

const playersKey = "players"

func getPlayers(req *http.Request) (int, error) {
	players := 5
	if plstrs, ok := req.Form[playersKey]; ok && len(plstrs) > 0 {
		pl, err := strconv.ParseInt(plstrs[0], 10, 32)
		if err != nil {
			return 0, err
		}
		players = int(pl)
	}
	if players < 2 {
		return players, errors.New(fmt.Sprintf("At least two players required, found %v", players))
	}
	return players, nil
}

func PlayHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintln(w, "<html><head><title>A game of Texas Hold'em</title>")
	fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(w, "</head><body><h1>A game of Texas Hold'em</h1>")
	players, err := getPlayers(req)
	if err != nil {
		// Use template to sanitise user input for security
		t := template.Must(template.New("error").Parse("<p>Could not parse players as integer: {{.}}</p></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	fmt.Fprintf(w, `<form method="get">Players: <input type="text" name="%v" value="%v"/><input type="submit" value="Rerun"/></form>`, playersKey, players)

	pack := poker.NewPack()
	pack.Shuffle()
	onTable, playerCards := pack.Deal(players)
	sortedOutcomes := poker.DealOutcomes(onTable, playerCards)
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
