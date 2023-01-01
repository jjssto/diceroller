package main

var rooms map[int]Room
var playerIds map[int][]int
var stats Statistic
var configuration ServerConfig

var MAX_TRIES_ID_GEN int = 100
var MAX_DICE int = 20

func main() {
	configuration.loadConfig("diceroller.conf")
	go cleanup(configuration)
	go runStatistics(configuration)
	initRand()

	rooms = make(map[int]Room)

	serve(configuration)
}
