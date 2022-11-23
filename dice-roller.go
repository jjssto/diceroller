package main

import (
	"github.com/gin-gonic/gin"
)

var rooms map[int]Room

func main() {
	initRand()
	rooms = make(map[int]Room)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", index)
	router.GET("/rolls/:id", getAllRolls)
	router.GET("/rolls/:id/:time", getRolls)
	router.POST("/roll", rollDice)
	router.POST("/", addRoomHandler)

	router.Run("localhost:8080")
}
