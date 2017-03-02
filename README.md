# Introduction

This project comprises some miscellaneous poker-related functions written for fun in [the Go programming language](http://golang.org/).

There is an HTTP front end, which so far provides the following features:

* "Play Holdem", which simulates a single hand of Texas Hold'em with a given number of players and displays the ranking of the hands
* "Simulate Holdem", which allows you to specify a number of known cards (both on the table and in your hand) and simulates a large number of hands of Texas Hold'em to see how likely various possible outcomes are. This gives an estimate of the conditional probabilities of the various game outcomes, given the cards that you know. (Poker strategy cannot be reduced to an algorithm purely based on these probabilities - you have to take your opponents' playing styles and betting behaviour into account, which is what makes poker an interesting game - but it is still very helpful to have a good sense of them.)
* "Starting Holdem cards", which compares the win probabilities from holding different starting pairs in Texas Hold'em. This information is useful when considering which hands to play and which to fold pre-flop. Simulations are done concurrently and inserted into the page in real-time using Angular.JS.
* "Play Omaha/8", which simulates a single hand of Omaha 8-or-better with a given number of players and displays the outcome.

# Installing and running locally

## Prerequisites

* You need to have [the Go compiler](http://golang.org/) installed on your computer (such that the ```go``` tool is available on your ```PATH```).
* You need a workspace set up according to [these instructions](https://golang.org/doc/code.html), for example:
  * ```mkdir $HOME/go```
  * ```export GOPATH=$HOME/go```
  * ```export PATH=$PATH:$GOPATH/bin```

## Running HTTP server

* Download this code into your workspace using ```go get github.com/amdw/gopoker```
* Build and install using ```go install github.com/amdw/gopoker```
* Run ```gopoker```.

## Running tests

To run the tests, you can run:

    go test github.com/amdw/gopoker/...

# Running on Heroku

The Gopoker HTTP server will run on the [Heroku](https://heroku.com) cloud platform using the standard Heroku buildpack for Go.

* Create a Heroku account and install the [Heroku toolbelt](https://devcenter.heroku.com/articles/heroku-command)
* Clone this repository
* Run ```heroku create -b https://github.com/heroku/heroku-buildpack-go``` to create the new Heroku application
* Run ```git push heroku master``` to push the code to Heroku
* Run ```heroku open``` to open the app in a web browser

# Licensing

This code was written as an exercise in Go programming and is not intended for any serious purposes whatsoever. For this reason, I have chosen to release it under a standard copyleft open source license which requires anyone who runs it on a network server to release any modifications they make under the same license.

This should cause no inconvenience to anyone who wants to use this software for educational purposes, while making it unattractive to anyone who wants to incorporate it into, for example, a commercial online poker site. (I have nothing in principle against online poker: I just prefer not to have others profit from incorporating my work into a proprietary product without compensating me.)

All files are Copyright 2013 Andrew Medworth and are released as free software under the GNU Affero GPL. See LICENSE file for more details.