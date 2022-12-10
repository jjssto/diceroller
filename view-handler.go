package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func viewHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"nbrCoC":       stats.nbrCoC,
		"nbrRezTech":   stats.nbrRezTech,
		"nbrGeneral":   stats.nbrGeneral,
		"nbrDiceRolls": stats.nbrDiceRolls,
		"nbrPlayers":   stats.nbrPlayer,
	})
}

func faviconHandler(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "favicon.png")
}

func viewCoC(c *gin.Context, id int) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "coc.html", gin.H{
		"title":     "Call of Cthulhu, Room #" + strconv.FormatInt(int64(id), 10),
		"player_id": session.Get("player_id"),
	})
}

func viewRezTech(c *gin.Context, id int) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "reztech.html", gin.H{
		"title":     "RezTech, Room #" + strconv.FormatInt(int64(id), 10),
		"player_id": session.Get("player_id"),
	})
}

func viewGeneral(c *gin.Context, id int) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "general.html", gin.H{
		"title":     "Room #" + strconv.FormatInt(int64(id), 10),
		"player_id": session.Get("player_id"),
	})
}

func viewGame(c *gin.Context) {
	roomId64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		displayError(c, err)
	}
	roomId := int(roomId64)
	r, ok := rooms[roomId]
	if !ok {
		displayError(c, nil)
		return
	}
	session := sessions.Default(c)
	playerId := session.Get("player_id")
	if playerId == nil {
		ok = false
	} else {
		_, ok = playerIds[playerId.(int)]
	}
	if !ok {
		playerId, ok := genPlayerId(roomId)
		if !ok {
			displayError(c, errors.New("error generationg player id"))
		}
		session.Set("player_id", playerId)
		session.Save()
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

func displayError(c *gin.Context, err interface{}) {
	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
		"title": "An error has occured",
	})
}

func viewError(c *gin.Context) {
	displayError(c, nil)
}
