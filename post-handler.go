package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func rollDice(c *gin.Context) {
	idStr := c.Param("id")
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 32)
	session := sessions.Default(c)
	player := session.Get("player_id").(int)
	col := json["color"]
	char := json["char"]
	mod, _ := strconv.ParseInt(json["mod"], 10, 32)
	action := json["action"]
	dice := json["dice"]
	diceArr := make([]int8, 0)
	if dice != "" {
	}
	r := rooms[int(id)]
	r.addPlayer(player, char, col)
	r.roll(diceArr, int(mod), int(player), action)
	rooms[int(id)] = r
	c.Status(200)
}

func addPlayer(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	room, _ := strconv.ParseInt(json["room"], 10, 32)
	id, _ := strconv.ParseInt(json["id"], 10, 32)
	name, _ := json["name"]
	color, _ := json["color"]
	r := rooms[int(room)]
	r.addPlayer(int(id), name, color)
	c.Status(http.StatusOK)
}

func addRoomHandler(c *gin.Context) {
	game, ok := c.GetPostForm("id")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	var g Game
	switch game {
	case "CoC":
		g = CoC
	case "RezTech":
		g = RezTech
	default:
		c.Status(http.StatusBadRequest)
		return
	}
	id, err := addRoom(g)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.AsciiJSON(http.StatusOK, id)
}
