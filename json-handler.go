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
