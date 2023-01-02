package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getColorOptions() string {
	str := "<option value=\"-\">-</option>\n"
	for key := range globConfig.Colors {
		str += fmt.Sprintf(
			"<option value=\"%s\">%s</option>\n",
			globConfig.Colors[key], key)
	}
	return str
}

func getTitle(room Room) string {
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

func checkOwnership(c *gin.Context, room Room) bool {
	playerId := getPlayerId(c, room.Id)
	if playerId == room.OwnerId {
		return true
	} else {
		return false
	}
}

func viewHome(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title":  globConfig.Title,
		"footer": globConfig.Footer,
	})
}

func viewCoC(c *gin.Context, id int) {
	room := globRooms[id]
	template.ParseFiles()
	c.HTML(http.StatusOK, "coc.html", gin.H{
		"title":       getTitle(room),
		"color":       room.Color,
		"room_id":     fmt.Sprintf("%d", room.Id),
		"is_owner":    fmt.Sprint(checkOwnership(c, room)),
		"result_cols": []string{"Roll"},
	})
}

func viewRezTech(c *gin.Context, id int) {
	room := globRooms[id]
	c.HTML(http.StatusOK, "reztech.html", gin.H{
		"title":       getTitle(room),
		"color":       room.Color,
		"room_id":     fmt.Sprintf("%d", room.Id),
		"is_owner":    fmt.Sprint(checkOwnership(c, room)),
		"result_cols": []string{"D12", "D8", "D6"},
	})
}

func viewGeneral(c *gin.Context, id int) {
	room := globRooms[id]
	c.HTML(http.StatusOK, "general.html", gin.H{
		"title":       getTitle(room),
		"color":       room.Color,
		"room_id":     fmt.Sprintf("%d", room.Id),
		"is_owner":    fmt.Sprint(checkOwnership(c, room)),
		"result_cols": []string{"D20", "D12", "D10", "D8", "D6", "D4"},
	})
}

func viewGame(c *gin.Context) {
	roomId64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		displayError(c, err)
	}
	roomId := int(roomId64)
	r, ok := globRooms[roomId]
	if !ok {
		displayError(c, nil)
		return
	}
	playerId := getPlayerId(c, roomId)
	if r.setOwnerId(playerId) {
		globRooms[roomId] = r
	}
	switch r.Game {
	case CoC:
		viewCoC(c, roomId)
	case RezTech:
		viewRezTech(c, roomId)
	case General:
		viewGeneral(c, roomId)
	default:
		displayError(c, err)
	}
}

func (room *Room) setOwnerId(playerId int) bool {
	if room.OwnerId == 0 {
		room.OwnerId = playerId
		return true
	} else {
		return false
	}
}

func getPlayerId(c *gin.Context, roomId int) int {
	session := sessions.Default(c)
	playerIdRaw := session.Get("player_id")
	var playerId int
	var ok bool
	if playerIdRaw == nil {
		ok = false
	} else {
		playerId = playerIdRaw.(int)
		_, ok = globPlayerIds[playerId]
	}
	if !ok {
		playerId, ok = genPlayerId(roomId)
		if !ok {
			displayError(c, errors.New("error generationg player id"))
		}
		session.Set("player_id", playerId)
		session.Save()
	}
	return playerId
}

func displayError(c *gin.Context, err interface{}) {
	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
		"title": "An error has occured",
	})
}

func viewError(c *gin.Context) {
	displayError(c, nil)
}

func viewStats(c *gin.Context) {
	c.HTML(http.StatusOK, "stats.html", gin.H{
		"nbrCoC":       globStats.nbrCoC,
		"nbrRezTech":   globStats.nbrRezTech,
		"nbrGeneral":   globStats.nbrGeneral,
		"nbrDiceRolls": globStats.nbrDiceRolls,
		"nbrPlayers":   globStats.nbrPlayer,
	})
}
