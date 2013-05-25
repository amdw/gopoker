package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"poker"
	"strconv"
	"strings"
)

func menu(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Poker</title></head><body><h1>Poker</h1>")
	fmt.Fprintf(w, `<ul><li><a href="/play">Play</a></li></ul>`)
	fmt.Fprintf(w, "</body></html>")
}

func formatCards(cards []poker.Card) string {
	cardStrings := make([]string, len(cards))
	for i, c := range cards {
		cardStrings[i] = c.HTML()
	}
	return strings.Join(cardStrings, ", ")
}

func playHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, "<html><head><title>A game of Texas Hold'em</title></head><body><h1>A game of Texas Hold'em</h1>")

	players := 5
	if plstrs, ok := req.Form["players"]; ok && len(plstrs) > 0 {
		pl, err := strconv.ParseInt(plstrs[0], 10, 32)
		if err != nil {
			fmt.Fprintf(w, "<p>Could not parse players as integer: %v</p></body></html>", err.Error())
			return
		}
		players = int(pl)
	}
	pack := poker.NewPack()
	pack.Shuffle()
	onTable, playerCards, sortedOutcomes := pack.PlayHoldem(players)
	fmt.Fprintf(w, "<h2>Table cards</h2><p>%v</p>", formatCards(onTable))
	fmt.Fprintf(w, "<h2>Player cards</h2><ul>")
	for player := 0; player < players; player++ {
		fmt.Fprintf(w, "<li>Player %v: %v</li>", player+1, formatCards(playerCards[player]))
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "<h2>Results</h2><table><tr><th>Position</th><th>Player</th><th>Hand</th><th>Cards</th></tr>")
	for i, outcome := range sortedOutcomes.Outcomes {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", i+1, outcome.Player, outcome.Level, formatCards(outcome.Cards))
	}
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</body></html>")
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	flag.Parse()

	log.Printf("Listening on port %v...\n", port)

	http.HandleFunc("/", menu)
	http.HandleFunc("/play", playHoldem)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
