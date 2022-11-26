package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getRollsHelper(c *gin.Context, all bool) {
	var roomId int
	roomId64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	} else {
		roomId = int(roomId64)
	}
	var rollNbr int
	if !all {
		rollNbr64, err := strconv.ParseInt(c.Param("roll_nbr"), 10, 64)
		if err != nil {
			rollNbr = 0
		} else {
			rollNbr = int(rollNbr64) + 1
		}
	} else {
		rollNbr = 0
	}
	_, ok := rooms[roomId]
	if !ok {
		c.JSON(0, "")
	}
	session := sessions.Default(c)
	playerId := session.Get("player_id").(int)
	data := "{"
	first := true
	for rollNbr < len(rooms[roomId].DiceRolls) {
		val := rooms[roomId].DiceRolls[rollNbr]
		if first {
			first = false
		} else {
			data += ","
		}
		isOwnRoll := 0
		if playerId == val.Player {
			isOwnRoll = 1
		}
		data += fmt.Sprintf("\"%d\": [%s, %d]", rollNbr, val.json(roomId), isOwnRoll)
		rollNbr++
	}
	data += "}"
	c.String(http.StatusOK, data)

}

func getAllRolls(c *gin.Context) {
	getRollsHelper(c, true)
}

func getRolls(c *gin.Context) {
	getRollsHelper(c, false)
}
