package Controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayersController struct {
}

func NewPlayersController() *PlayersController {
	controller := PlayersController{}
	return &controller
}

func (controller *PlayersController) Route(router *gin.Engine) {
	group := router.Group("players")
	group.GET("/players/", controller.homeHandler)
}

func (controller *PlayersController) homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}
