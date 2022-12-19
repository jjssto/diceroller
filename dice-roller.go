package main

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var rooms map[int]Room
var stats Statistic

var playerIds map[int][]int
var MAX_TRIES_ID_GEN int = 100
var INACTIVE_DELETE_DELAY string = "4h"
var MAX_DICE int = 200

func cleanup() {
	for {
		time.Sleep(15 * time.Minute)
		deleteOldGames(rooms, playerIds)
	}
}

func runStatistics() {
	for {
		time.Sleep(time.Minute)
		stats, _ = updateStatistics(rooms, playerIds)
	}
}

func main() {
	go cleanup()
	go runStatistics()
	initRand()
	rooms = make(map[int]Room)
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.LoadHTMLGlob("templates/*")
	router.Static("/js", "./js")
	router.Static("/css", "./css")
	router.Static("/rec", "./rec")
	router.StaticFile("/favicon.png", "./favicon.png")
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", viewHome)
	router.GET("/stats", viewStats)
	router.GET("/error", viewError)
	router.GET("/room/:id", viewGame)
	router.GET("/rolls/:id", getAllRolls)
	router.GET("/rolls/:id/:roll_nbr", getRolls)
	router.POST("/room/:id", rollDice)
	router.POST("/", addRoomHandler)
	router.Run("localhost:9000")

}
