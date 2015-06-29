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
package main

import (
	"flag"
	"fmt"
	"github.com/amdw/gopoker/poker_http"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Try some sensible defaults to get the HTML path
func defaultHtmlBaseDir() string {
	// First try something GOPATH-relative as this is most likely to work
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		gopaths := filepath.SplitList(gopath)
		for _, p := range gopaths {
			candidate := path.Join(p, "src", "github.com", "amdw", "gopoker", "html")
			log.Println("Trying", candidate)
			dirInfo, err := os.Stat(candidate)
			if err == nil && dirInfo.IsDir() {
				return candidate
			}
		}
	}
	// If that didn't work, try the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Warning: could not get current working directory:", err)
		return "./html"
	}
	return path.Join(currentDir, "html")
}

func main() {
	var port int
	var htmlBaseDir string
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	flag.StringVar(&htmlBaseDir, "htmlbasedir", defaultHtmlBaseDir(), "Base directory containing HTML")
	flag.Parse()

	dirInfo, err := os.Stat(htmlBaseDir)
	if err != nil {
		log.Fatalf("Could not stat HTML base dir '%v' - provide it on the command line or ensure GOPATH is set correctly.\n(Error: %v)", htmlBaseDir, err)
	}
	if !dirInfo.IsDir() {
		log.Fatalf("htmlbasedir '%v' is not a directory", htmlBaseDir)
	}

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
