package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func viewHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Dice roller",
	})
}

func viewCoC(c *gin.Context, id int64) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "coc.html", gin.H{
		"title":     "Call of Cathulu, Room #" + strconv.FormatInt(id, 10),
		"player_id": session.Get("player_id"),
	})
}

func viewGame(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	r, ok := rooms[int(id)]
	session := sessions.Default(c)
	x := session.Get("player_id")
	if x == nil {
		player := rand.Intn(99999998) + 1
		session.Set("player_id", player)
		session.Save()
	}
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	switch r.Game {
	case CoC:
		viewCoC(c, id)
	case RezTech:
		c.Redirect(http.StatusTemporaryRedirect, "/")
	default:
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
