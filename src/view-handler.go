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
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var a []string = []string{
	"Attack", "Defend", "Dodge", "Craft", "Social", "Athletics",
	"Magic", "Perception",
}

type Input struct {
	Id      string
	Label   string
	Options []int
}

func (room *Room) getTitle() string {
	var ret string
	if len(room.Name) > 0 {
		ret = room.Name
	} else {
		ret = fmt.Sprintf("Room #%d", room.Id)
	}
	switch room.Game {
	case CoC:
		ret += " [Call of Cthulhu]"
	case RezTech:
		ret += " [RezTech]"
	case General:
		ret += ""
	}
	return ret
}

func viewHome(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title": globConfig.Title,
	})
}

func viewCoC(c *gin.Context, room Room) {
	s := make([]int, 10)
	for i := range s {
		s[i] = i + 1
	}
	c.HTML(http.StatusOK, "coc.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"D100", "D10"},
		"d20":         Input{Id: "s_d20", Label: "D20:", Options: s},
		"d12":         Input{Id: "s_d12", Label: "D12:", Options: s},
		"d10":         Input{Id: "s_d10", Label: "D10:", Options: s},
		"d00":         Input{Id: "s_d00", Label: "D00:", Options: s},
		"d8":          Input{Id: "s_d8", Label: "D8:", Options: s},
		"d6":          Input{Id: "s_d6", Label: "D6:", Options: s},
		"d4":          Input{Id: "s_d4", Label: "D4:", Options: s},
		"actions":     a,
	})
}

func viewRezTech(c *gin.Context, room Room) {
	s := make([]int, 10)
	for i := range s {
		s[i] = i + 1
	}
	c.HTML(http.StatusOK, "reztech.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"D12", "D8", "D6"},
		"attr_opt":    Input{Id: "s_attribute", Label: "Attribute:", Options: s},
		"skill_opt":   Input{Id: "s_skill", Label: "Skill:", Options: s},
		"d20":         Input{Id: "s_d20", Label: "D20:", Options: s},
		"d12":         Input{Id: "s_d12", Label: "D12:", Options: s},
		"d10":         Input{Id: "s_d10", Label: "D10:", Options: s},
		"d00":         Input{Id: "s_d00", Label: "D00:", Options: s},
		"d8":          Input{Id: "s_d8", Label: "D8:", Options: s},
		"d6":          Input{Id: "s_d6", Label: "D6:", Options: s},
		"d4":          Input{Id: "s_d4", Label: "D4:", Options: s},
		"actions":     a,
	})
}

func viewGeneral(c *gin.Context, room Room) {
	s := make([]int, 10)
	for i := range s {
		s[i] = i + 1
	}
	c.HTML(http.StatusOK, "general.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"D20", "D12", "D10", "D8", "D6", "D4"},
		"d20":         Input{Id: "s_d20", Label: "D20:", Options: s},
		"d12":         Input{Id: "s_d12", Label: "D12:", Options: s},
		"d10":         Input{Id: "s_d10", Label: "D10:", Options: s},
		"d00":         Input{Id: "s_d00", Label: "D00:", Options: s},
		"d8":          Input{Id: "s_d8", Label: "D8:", Options: s},
		"d6":          Input{Id: "s_d6", Label: "D6:", Options: s},
		"d4":          Input{Id: "s_d4", Label: "D4:", Options: s},
		"actions":     a,
	})
}
func viewGame(c *gin.Context) {
	roomId64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		displayError(c, err)
		return
	}
	roomId := int(roomId64)

	db := DB{Configured: false}
	db.connect(false)
	userToken, _ := db.getToken(c)
	room, err := db.getRoom(roomId, userToken)
	db.close()
	if err != nil {
		displayError(c, err)
		return
	}
	switch room.Game {
	case CoC:
		viewCoC(c, room)
	case RezTech:
		viewRezTech(c, room)
	case General:
		viewGeneral(c, room)
	default:
		displayError(c, err)
	}
}

func displayError(c *gin.Context, err interface{}) {
	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
		"title": "An error has occured",
	})
}

func viewError(c *gin.Context) {
	displayError(c, nil)
}

func viewRooms(c *gin.Context) {
	db := DB{}
	db.connect(false)
	userToken, _ := db.getToken(c)
	ownRooms, err1 := db.getOwnRooms(userToken)
	allRooms, err2 := db.getAllRooms(userToken)
	db.close()
	if err1 != nil || err2 != nil {
		displayError(c, nil)
		return
	}
	c.HTML(http.StatusOK, "rooms.html", gin.H{
		"title":     "Rooms",
		"own_rooms": ownRooms,
		"all_rooms": allRooms,
	})
}
