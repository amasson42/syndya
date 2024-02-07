package Controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PlayersController struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewPlayersController() *PlayersController {
	controller := PlayersController{}
	return &controller
}

func (controller *PlayersController) Route(router *gin.Engine) {
	router.GET("search", controller.searchGame)
	router.GET("players", controller.getPlayers)
}

func (controller *PlayersController) searchGame(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		log.Printf("Received: %s\n", message)

		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, socket!"))

		if err != nil {
			log.Println("Error writing message: ", err)
			break
		}

	}
}

func (controller *PlayersController) getPlayers(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]any{
		"001": map[string]any{
			"rating": 100,
			"name":   "mike",
		},
		"002": map[string]any{
			"rating": 200,
			"name":   "joe",
		},
		"003": map[string]any{
			"rating": 300,
			"name":   "roger",
		},
	})
}
