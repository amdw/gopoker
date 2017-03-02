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

// Try some sensible defaults to get the static content path
func defaultStaticBaseDir() string {
	// First try something GOPATH-relative as this is most likely to work
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		gopaths := filepath.SplitList(gopath)
		for _, p := range gopaths {
			candidate := path.Join(p, "src", "github.com", "amdw", "gopoker", "static")
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
		return "./static"
	}
	return path.Join(currentDir, "static")
}

func main() {
	var port int
	var staticBaseDir string
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	flag.StringVar(&staticBaseDir, "staticbasedir", defaultStaticBaseDir(), "Base directory containing static content")
	flag.Parse()

	dirInfo, err := os.Stat(staticBaseDir)
	if err != nil {
		log.Fatalf("Could not stat static content base dir '%v' - provide it on the command line or ensure GOPATH is set correctly.\n(Error: %v)", staticBaseDir, err)
	}
	if !dirInfo.IsDir() {
		log.Fatalf("staticbasedir '%v' is not a directory", staticBaseDir)
	}

	log.Println("Using static content base dir", staticBaseDir)
	log.Printf("Listening on port %v...\n", port)

	http.HandleFunc("/", poker_http.Menu)
	http.HandleFunc("/holdem/play", poker_http.PlayHoldem)
	http.HandleFunc("/holdem/simulate", poker_http.SimulateHoldem(staticBaseDir))
	http.HandleFunc("/holdem/startingcards", poker_http.StartingCards(staticBaseDir))
	http.HandleFunc("/holdem/startingcards/sim", poker_http.SimulateStartingCards)
	http.HandleFunc("/omaha8/play", poker_http.PlayOmaha8)
	http.HandleFunc("/omaha8/simulate", poker_http.SimulateOmaha8)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
