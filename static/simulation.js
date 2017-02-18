var app = angular.module('simulatorApp', ['ngSanitize']);

app.controller('simulatorController', function($scope, $window, $sce) {
    $scope.playerCount = initPlayerCount;
    $scope.yourCards = initYourCards;
    $scope.tableCards = initTableCards;
    $scope.simulationCount = initSimCount;

    var legalRanks = ["A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"];
    var legalSuits = ["C", "D", "H", "S"];

    $scope.fewerPlayers = function() {
        $scope.playerCount = Math.max(2, $scope.playerCount - 1);
    };
    $scope.morePlayers = function() {
        $scope.playerCount += 1;
    };

    var displayCard = function(card) {
        var rank = card.toUpperCase().slice(0, -1);
        var suit = card.toUpperCase().slice(-1);
        var displaySuit;
        switch (suit) {
            case "C":
                displaySuit = "&#9827;";
                break;
            case "D":
                displaySuit = '<span style="color:red">&#9830;</span>';
                break;
            case "H":
                displaySuit = '<span style="color:red">&#9829;</span>';
                break;
            case "S":
                displaySuit = "&#9824;";
                break;
            default:
                displaySuit = "?";
        }
        // This is trusted HTML so don't allow arbitrary stuff through
        var displayRank = legalRanks.indexOf(rank) >= 0 ? rank : "?";
        return displayRank + displaySuit;
    };
    var displayCards = function(cards) {
        return cards.map(displayCard).join(", ");
    };
    var cardsUri = function(cards) {
        return encodeURIComponent(cards.join(","));
    };

    $scope.displayYourCards = function() {
        return $sce.trustAsHtml(displayCards($scope.yourCards));
    };
    $scope.yourCardsUri = function() {
        return cardsUri($scope.yourCards);
    };
    $scope.displayTableCards = function() {
        return $sce.trustAsHtml(displayCards($scope.tableCards));
    };
    $scope.tableCardsUri = function() {
        return cardsUri($scope.tableCards);
    };

    var remainingPack = function() {
        var result = [];
        for (i = 0; i < legalSuits.length; i++) {
            for (j = 0; j < legalRanks.length; j++) {
                var card = legalRanks[j] + legalSuits[i];
                if ($scope.yourCards.indexOf(card) < 0 && $scope.tableCards.indexOf(card) < 0) {
                    result.push(card);
                }
            }
        }
        return result;
    };
    var randomRemainingCards = function(count) {
        var pack = remainingPack();
        var result = [];
        while (result.length < count && pack.length > 0) {
            var idx = Math.floor(Math.random() * pack.length);
            result.push(pack[idx]);
            pack.splice(idx, 1);
        }
        return result;
    };
    $scope.yourCardsRandom = function() {
        var newCards = randomRemainingCards(Math.max(0, 2 - $scope.yourCards.length));
        $scope.yourCards = $scope.yourCards.concat(newCards);
    };
    var tableCardsRandom = function(desiredNumber) {
        var newCards = randomRemainingCards(Math.max(0, desiredNumber - $scope.tableCards.length));
        $scope.tableCards = $scope.tableCards.concat(newCards);
    };
    $scope.flopRandom = function() {
        tableCardsRandom(3);
    };
    $scope.turnRandom = function() {
        tableCardsRandom(4);
    };
    $scope.riverRandom = function() {
        tableCardsRandom(5);
    };

    $scope.rerun = function() {
        var parts = ["players=" + $scope.playerCount];
        if ($scope.yourCards.length > 0) {
            parts.push("yours=" + $scope.yourCardsUri());
        }
        if ($scope.tableCards.length > 0) {
            parts.push("table=" + $scope.tableCardsUri());
        }
        parts.push("simcount=" + $scope.simulationCount);
        $window.location.href = "/simulate?" + parts.join("&");
    };
});
