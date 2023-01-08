// This code is licensed under the MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	userToken, err := getToken(c)
	if err != nil {
		userToken = 0
	} else if userToken == 0 {
		userToken, _ = db.getToken(c)
	}

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
	userToken, err := getToken(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	} else if userToken == 0 {
		_, _ = db.getToken(c)
	}
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
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userToken, err := getToken(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
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

func deleteRoomHandler(c *gin.Context) {
	var json map[string]string
	userToken, err := getToken(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	if err = c.BindJSON(&json); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	roomId, _ := strconv.ParseInt(json["room_id"], 10, 32)
	db := DB{}
	db.connect(false)
	nbr, _, err := db.deleteRoom(userToken, roomId)
	db.close()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	if nbr > 0 {
		c.String(http.StatusOK, fmt.Sprint(nbr))
	} else {
		c.Status(http.StatusForbidden)
	}
}
