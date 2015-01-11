package main

import (
	"flag"
	"fmt"
	"gopoker/poker_http"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	flag.Parse()

	log.Printf("Listening on port %v...\n", port)

	http.HandleFunc("/", poker_http.Menu)
	http.HandleFunc("/play", poker_http.PlayHoldem)
	http.HandleFunc("/simulate", poker_http.SimulateHoldem)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
