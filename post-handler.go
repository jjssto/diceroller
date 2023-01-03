package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func checkDiceArr(diceArr []int8) ([]int8, bool) {
	if len(diceArr) > 200 {
		ret := make([]int8, 200)
		for i := 0; i < MAX_DICE; i++ {
			ret[i] = diceArr[i]
		}
		return ret, false
	} else {
		return diceArr, true
	}
}

func rollDice(c *gin.Context) {
	idStr := c.Param("id")
	var data map[string]string
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 32)
	session := sessions.Default(c)
	userToken := session.Get("player_id").(int)
	col := data["color"]
	char := data["char"]
	mod, _ := strconv.ParseInt(data["mod"], 10, 32)
	action := data["action"]
	dice := data["dice"]
	diceArr := make([]int8, 0)
	if dice != "" {
		json.Unmarshal([]byte(dice), &diceArr)
	}
	diceArr, _ = checkDiceArr(diceArr)
	db := DB{Configured: false}
	db.connect(false)
	err = db.roll(
		int(id), userToken, char, col, action, int(mod), diceArr,
	)
	db.close()
	if err != nil {
		c.Status(http.StatusForbidden)
	} else {
		c.Status(http.StatusOK)
	}
}

func addRoomHandler(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	game := json["id"]
	var g Game
	switch game {
	case "CoC":
		g = CoC
	case "RezTech":
		g = RezTech
	default:
		g = General
	}
	db := DB{Configured: false}
	db.connect(false)
	id, err := db.createRoom(g)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	db.close()
	c.AsciiJSON(http.StatusOK, id)
}

func changeRoomSettings(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	session := sessions.Default(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userToken := session.Get("player_id").(int)
	roomId, _ := strconv.ParseInt(json["room_id"], 10, 32)
	roomName := json["room_name"]
	color := json["color"]

	db := DB{Configured: false}
	if err := db.connect(false); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	err = db.changeRoomSettings(int(roomId), userToken, roomName, color)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	c.Status(http.StatusOK)
}
