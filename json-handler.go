package main

import (
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
			rollNbr = int(rollNbr64)
		}
	} else {
		rollNbr = 0
	}
	session := sessions.Default(c)
	userToken := session.Get("player_id").(int)
	//loc := getTimeZone(c)

	db := DB{Configured: false}
	db.connect(false)
	data := db.getRolls(roomId, userToken, rollNbr)
	db.close()
	c.Writer.Flush()
	c.Writer.WriteString(data)
	c.Status(http.StatusOK)

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

type Color struct {
	Text string
	Code string
}
