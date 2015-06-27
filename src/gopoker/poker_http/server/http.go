package main

import (
	"flag"
	"fmt"
	"gopoker/poker_http"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	var port int
	var htmlBaseDir string
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Warning: could not get current working directory:", err)
		currentDir = "./"
	}
	flag.StringVar(&htmlBaseDir, "htmlbasedir", path.Join(currentDir, "html"), "Base directory containing HTML")
	flag.Parse()

	log.Println("Using HTML base dir", htmlBaseDir)
	log.Printf("Listening on port %v...\n", port)

	http.HandleFunc("/", poker_http.Menu)
	http.HandleFunc("/play", poker_http.PlayHoldem)
	http.HandleFunc("/simulate", poker_http.SimulateHoldem)
	http.HandleFunc("/startingcards", poker_http.StartingCards(htmlBaseDir))
	http.HandleFunc("/startingcards/sim", poker_http.SimulateStartingCards)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
