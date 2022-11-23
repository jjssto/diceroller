package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func rollDice(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		panic(err)
	}
	id, _ := strconv.ParseInt(json["id"], 10, 32)
	player, _ := strconv.ParseInt(json["player"], 10, 32)
	mod, _ := strconv.ParseInt(json["mod"], 10, 32)
	action, _ := json["action"]
	dice, _ := json["dice"]
	diceArr := make([]int8, 0)
	if dice != "" {
	}
	r := rooms[int(id)]
	r.roll(diceArr, int(mod), int(player), action)
	c.Status(200)
}

func addPlayer(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		panic(err)
	}
	room, _ := strconv.ParseInt(json["room"], 10, 32)
	id, _ := strconv.ParseInt(json["id"], 10, 32)
	name, _ := json["name"]
	color, _ := json["color"]
	r := rooms[int(room)]
	r.addPlayer(int(id), name, color)
	c.Status(200)
}

func addRoomHandler(c *gin.Context) {
	game, ok := c.GetPostForm("id")
	if !ok {

	}
	//var json map[string]string
	//err := c.BindJSON(&json)
	//if err != nil {
	//	panic(err)
	//}
	//game, _ := json["game"]
	var g Game
	switch game {
	case "CoC":
		g = CoC
	case "RezTech":
		g = RezTech
	default:
		c.Status(401)
	}
	id, err := addRoom(g)
	if err != nil {

	}
	c.Status(id)
	//c.Redirect(fmt.Sprintf("/%d", id))
}
