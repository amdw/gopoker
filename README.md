# Introduction

This project comprises some miscellaneous poker-related functions written for fun in [the Go programming language](http://golang.org/).

So far, there are the following features:

* a command-line program play_poker.go which randomly generates a large number of five-card poker hands and displays the frequencies of the various hand types
* a program poker_http.go which starts a HTTP server giving a simple interface to the following:
  * "Play", which simulates a single hand of Texas Hold'em with a given number of players and displays the ranking of the hands
  * "Simulate", which allows you to specify a number of known cards (both on the table and in your hand) and simulates a large number of hands of Texas Hold'em to see how likely various possible outcomes are. This is perhaps the most interesting feature, as it gives an estimate of the conditional probabilities of the various game outcomes, given the cards that you know.

# Installing and Running

An easy way to download and run this code is to clone the repository, add the resulting gopoker directory to your GOPATH, cd into it, and run

    go install gopoker/...

This should build and install both programs in the bin directory. You can then run either ```bin/play_poker``` or ```bin/poker_http```.

Obviously, this requires you to have Go installed on your computer.

# Licensing

This code was written as an exercise in Go programming and is not intended for any serious purposes whatsoever. For this reason, I have chosen to release it under a standard copyleft open source license which requires anyone who runs it on a network server to release any modifications they make under the same license.

This should cause no inconvenience to anyone who wants to use this software for educational purposes, while making it unattractive to anyone who wants to incorporate it into, for example, a commercial online poker site.

All files are Copyright 2013 Andrew Medworth and are released as free software under the GNU Affero GPL. See LICENSE file for more details.