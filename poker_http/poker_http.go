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
	"errors"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"net/http"
	"strconv"
	"strings"
)

func Menu(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "<html><head>")
	fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(w, "<title>Poker</title></head><body><h1>Poker</h1><ul>")
	fmt.Fprintln(w, "<li>Texas Holdem<ul>")
	fmt.Fprintln(w, `<li><a href="/holdem/play">Play</a></li>`)
	fmt.Fprintln(w, `<li><a href="/holdem/simulate">Simulate</a></li>`)
	fmt.Fprintln(w, `<li><a href="/holdem/startingcards">Starting cards</a></li>`)
	fmt.Fprintln(w, "</ul></li>")
	fmt.Fprintln(w, "<li>Omaha/8<ul>")
	fmt.Fprintln(w, `<li><a href="/omaha8/play">Play</a></li>`)
	fmt.Fprintln(w, "</ul></li>")
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
