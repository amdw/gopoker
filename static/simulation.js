var app = angular.module('simulatorApp', ['ngSanitize']);

app.controller('simulatorController', function($scope, $window, $sce) {
    $scope.playerCount = initPlayerCount;
    $scope.yourCards = initYourCards;
    $scope.tableCards = initTableCards;
    $scope.simulationCount = initSimCount;
    $scope.potSize = 1000;

    $scope.legalRanks = ["2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"];
    $scope.legalSuits = ["C", "D", "H", "S"];

    $scope.yourPendingSuit = "";
    $scope.yourPendingRank = "";
    $scope.tablePendingSuit = "";
    $scope.tablePendingRank = "";

    $scope.fewerPlayers = function() {
        $scope.playerCount = Math.max(2, $scope.playerCount - 1);
    };
    $scope.morePlayers = function() {
        $scope.playerCount += 1;
    };

    $scope.displaySuit = function(suit) {
        switch (suit) {
            case "C":
                return $sce.trustAsHtml("&#9827;");
            case "D":
                return $sce.trustAsHtml('<span style="color:red">&#9830;</span>');
            case "H":
                return $sce.trustAsHtml('<span style="color:red">&#9829;</span>');
            case "S":
                return $sce.trustAsHtml("&#9824;");
        }
        return "?";
    };
    var displayCard = function(card) {
        var rank = card.toUpperCase().slice(0, -1);
        var suit = card.toUpperCase().slice(-1);
        // This is trusted HTML so don't allow arbitrary stuff through
        var displayRank = $scope.legalRanks.indexOf(rank) >= 0 ? rank : "?";
        return displayRank + $scope.displaySuit(suit);
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
        for (i = 0; i < $scope.legalSuits.length; i++) {
            for (j = 0; j < $scope.legalRanks.length; j++) {
                var card = $scope.legalRanks[j] + $scope.legalSuits[i];
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

    $scope.deleteOneYourCard = function() {
        $scope.yourCards.splice(-1, 1);
    };
    $scope.deleteOneTableCard = function() {
        $scope.tableCards.splice(-1, 1);
    };

    var addYourPendingCard = function() {
        if ($scope.legalSuits.indexOf($scope.yourPendingSuit) >= 0 && $scope.legalRanks.indexOf($scope.yourPendingRank) >= 0) {
            var card = $scope.yourPendingRank + $scope.yourPendingSuit;
            if ($scope.yourCards.indexOf(card) < 0 && $scope.tableCards.indexOf(card) < 0) {
                $scope.yourCards.push(card);
            }
            $scope.yourPendingRank = "";
            $scope.yourPendingSuit = "";
        }
    };
    var addTablePendingCard = function() {
        if ($scope.legalSuits.indexOf($scope.tablePendingSuit) >= 0 && $scope.legalRanks.indexOf($scope.tablePendingRank) >= 0) {
            var card = $scope.tablePendingRank + $scope.tablePendingSuit;
            if ($scope.tableCards.indexOf(card) < 0 && $scope.yourCards.indexOf(card) < 0) {
                $scope.tableCards.push(card);
            }
            $scope.tablePendingRank = "";
            $scope.tablePendingSuit = "";
        }
    }

    $scope.setYourPendingSuit = function(suit) {
        $scope.yourPendingSuit = suit;
        addYourPendingCard();
    };
    $scope.setYourPendingRank = function(rank) {
        $scope.yourPendingRank = rank;
        addYourPendingCard();
    };
    $scope.setTablePendingSuit = function(suit) {
        $scope.tablePendingSuit = suit;
        addTablePendingCard();
    };
    $scope.setTablePendingRank = function(rank) {
        $scope.tablePendingRank = rank;
        addTablePendingCard();
    };

    $scope.yourCardsEmpty = function() {
        return $scope.yourCards.length < 1;
    };
    $scope.yourCardsFull = function() {
        return $scope.yourCards.length >= 2;
    };
    $scope.tableCardsEmpty = function() {
        return $scope.tableCards.length < 1;
    };
    $scope.tableCardsFull = function() {
        return $scope.tableCards.length >= 5;
    };

    $scope.yourSuitButtonClasses = function(suit) {
        return {'active': suit == $scope.yourPendingSuit && !$scope.yourCardsFull()};
    };
    $scope.yourRankButtonClasses = function(rank) {
        return {'active': rank == $scope.yourPendingRank && !$scope.yourCardsFull()};
    };
    $scope.tableSuitButtonClasses = function(suit) {
        return {'active': suit == $scope.tablePendingSuit && !$scope.tableCardsFull()};
    };
    $scope.tableRankButtonClasses = function(rank) {
        return {'active': rank == $scope.tablePendingRank && !$scope.tableCardsFull()};
    };

    $scope.potOddsMessage = function() {
        if (potOddsBreakEven == Infinity) {
            return $sce.trustAsHtml("<b>Any</b> bet size has positive expected value! :)");
        } else {
            var maxBetSize = Math.floor(parseInt($scope.potSize) * potOddsBreakEven);
            return $sce.trustAsHtml("A bet of up to <b>" + maxBetSize.toLocaleString() + "</b>" +
                                    " (" + Math.round(potOddsBreakEven * 100) + "% of the pot) " +
                                    "has positive expected value.");
        }
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
        parts.push("runsim=true");
        $window.location.href = "/simulate?" + parts.join("&");
    };
});
