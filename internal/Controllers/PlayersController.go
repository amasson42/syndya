package Controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PlayersController struct {
	upgrader    websocket.Upgrader
	playersBank Models.SearchingPlayersBank
	openSockets map[int]*websocket.Conn
}

func NewPlayersController(playersBank Models.SearchingPlayersBank) *PlayersController {
	controller := PlayersController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		playersBank: playersBank,
		openSockets: map[int]*websocket.Conn{},
	}
	return &controller
}

func (controller *PlayersController) Route(router *gin.Engine) {
	router.GET("search", controller.searchGame)
	router.GET("players", controller.getPlayers)
}

func (controller *PlayersController) searchGame(c *gin.Context) {
	conn, err := controller.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading request: ", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	playerId := controller.playersBank.CreateSearchingPlayer()

	controller.openSockets[playerId] = conn

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Closed connection with %v", playerId)
		return nil
	})
	defer conn.Close()
	defer delete(controller.openSockets, playerId)
	defer controller.playersBank.DeleteSearchingPlayer(playerId)

	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("id %v", playerId)))
	if err != nil {
		log.Println("Error writing message: ", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		cmd := strings.SplitN(string(message[:]), " ", 2)

		if len(cmd) >= 1 {
			switch cmd[0] {
			case "meta":
				if len(cmd) >= 2 {
					parameters := strings.SplitN(cmd[1], "=", 2)
					if len(parameters) == 2 {
						metaKey := parameters[0]
						metaValue := parameters[1]
						controller.playersBank.UpdateSearchingPlayerMetadata(playerId, metaKey, metaValue)
					}
				}
			}
		}

		// err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, socket!"))
		// if err != nil {
		// 	log.Println("Error writing message: ", err)
		// 	break
		// }

	}
}

func (controller *PlayersController) getPlayers(c *gin.Context) {
	allPlayers := controller.playersBank.GetAllSearchingPlayers()
	c.JSON(http.StatusOK, allPlayers)
}
