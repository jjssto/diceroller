package main

var globRooms map[int]Room
var globUserIds map[int][]int
var globStats Statistic
var globConfig ServerConfig

var MAX_TRIES_ID_GEN int = 100
var MAX_DICE int = 20

func main() {
	globConfig.loadConfig("diceroller.conf")
	go cleanup(globConfig)
	go runStatistics(globConfig)
	initRand()

	globRooms = make(map[int]Room)

	serve(globConfig)
}
