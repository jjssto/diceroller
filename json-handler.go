package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	loc := getTimeZone(c)
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
		data += fmt.Sprintf("\"%d\": [%s, %d]", rollNbr, val.json(roomId, loc), isOwnRoll)
		rollNbr++
	}
	data += "}"
	c.String(http.StatusOK, data)

}

func getTimeZone(c *gin.Context) *time.Location {
	offsetStr := c.Request.Header.Values("ts_offset")
	offset := 0
	if len(offsetStr) > 0 {
		offsetInt64, err := strconv.ParseInt(offsetStr[0], 10, 32)
		if err == nil {
			offset = int(offsetInt64)
		}
	}
	return time.FixedZone("", offset)
}

func getAllRolls(c *gin.Context) {
	getRollsHelper(c, true)
}

func getRolls(c *gin.Context) {
	getRollsHelper(c, false)
}
