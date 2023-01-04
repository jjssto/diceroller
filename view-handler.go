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

func checkOwnership(c *gin.Context, room Room) bool {
	return false
}

func viewHome(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title": globConfig.Title,
	})
}

func viewCoC(c *gin.Context, room Room) {
	template.ParseFiles()
	c.HTML(http.StatusOK, "coc.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"Roll"},
	})
}

func viewRezTech(c *gin.Context, room Room) {
	c.HTML(http.StatusOK, "reztech.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"D12", "D8", "D6"},
	})
}

func viewGeneral(c *gin.Context, room Room) {
	c.HTML(http.StatusOK, "general.html", gin.H{
		"title":       room.getTitle(),
		"color":       room.Color,
		"room_id":     room.Id,
		"is_owner":    room.IsOwner,
		"result_cols": []string{"D20", "D12", "D10", "D8", "D6", "D4"},
	})
}

func viewGame(c *gin.Context) {
	roomId64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		displayError(c, err)
		return
	}
	roomId := int(roomId64)
	session := sessions.Default(c)
	oldToken := session.Get("player_id").(int)

	db := DB{Configured: false}
	db.connect(false)
	userToken, _, err := db.createToken(oldToken)
	if err != nil {
		displayError(c, err)
		return
	}
	if userToken != oldToken {
		session.Set("player_id", userToken)
		session.Save()
	}
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

func getPlayerId(c *gin.Context, roomId int) int {
	session := sessions.Default(c)
	playerIdRaw := session.Get("player_id")
	var playerId int
	var ok bool
	if playerIdRaw == nil {
		ok = false
	} else {
		playerId = playerIdRaw.(int)
		_, ok = globUserIds[playerId]
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
