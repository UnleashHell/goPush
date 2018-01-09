package controllers

import (
	"github.com/gin-gonic/gin"
	"goPush/lib/push"
	"net/http"
)

type PushController struct{}

const (
	NO_ERROR    = 0
	PARAM_ERROR = 1
	SERVER_BUSY = 2
)

type params struct {
	Token string `form:"token" binding:"required"`
	Alert string `form:"alert" binding:"required"`
	Badge int    `form:"badge"`
	Sound string `form:"sound"`
}

var message = push.Message{}

func (this *PushController) Push(g *gin.Context) {
	var form params
	if g.Bind(&form) == nil {
		message := message.CreateMessage(form.Token, form.Alert, form.Sound, form.Badge)
		result := push.IosInstance.Push(message)
		if result {
			g.JSON(http.StatusOK, gin.H{"error": NO_ERROR, "msg": "success"})

		} else {
			g.JSON(http.StatusOK, gin.H{"error": SERVER_BUSY, "msg": "server busy"})

		}
	} else {
		g.JSON(http.StatusOK, gin.H{"error": PARAM_ERROR, "msg": "param error"})
	}

}
