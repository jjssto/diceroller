package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAllRolls(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		panic(err)
	}
	_, ok := rooms[int(id)]
	if !ok {
		c.JSON(0, "")
	}
	c.JSON(0, rooms[int(id)].DiceRolls)
}

func getRolls(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	time, err := strconv.ParseInt(c.Param("time"), 10, 64)
	if err != nil {
		panic(err)
	}
	_, ok := rooms[int(id)]
	if !ok {
		c.JSON(0, "")
	}
	diceRolls := make(map[int64]DiceRoll)
	r := rooms[int(id)].DiceRolls
	for idx, val := range r {
		if idx > time {
			diceRolls[idx] = val
		}
	}
	c.JSON(0, diceRolls)
}
