package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var rooms map[int]Room

func main() {
	initRand()
	rooms = make(map[int]Room)
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.LoadHTMLGlob("templates/*")
	router.Static("/js", "./js")
	router.Static("/css", "./css")
	router.Static("/rec", "./rec")
	router.GET("/", viewHome)
	router.GET("/room/:id", viewGame)
	router.GET("/rolls/:id", getAllRolls)
	router.GET("/rolls/:id/:time", getRolls)
	router.POST("/room/:id", rollDice)
	router.POST("/", addRoomHandler)

	router.Run("localhost:8080")
}
