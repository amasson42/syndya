package Controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PlayersController struct {
	upgrader         websocket.Upgrader
	playersBank      Models.SearchingPlayersBank
	connections      map[int]*websocket.Conn
	connectionsMutex sync.Mutex
}

func NewPlayersController(playersBank Models.SearchingPlayersBank) *PlayersController {
	controller := PlayersController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		playersBank: playersBank,
		connections: map[int]*websocket.Conn{},
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

	controller.connectionsMutex.Lock()
	controller.connections[playerId] = conn
	controller.connectionsMutex.Unlock()

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Closed connection with %v", playerId)
		return nil
	})

	defer func() {
		controller.playersBank.DeleteSearchingPlayer(playerId)
		controller.connectionsMutex.Lock()
		delete(controller.connections, playerId)
		controller.connectionsMutex.Unlock()
		conn.Close()
	}()

	controller.sendPlayerId(playerId)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		controller.interpretWebSocketMessage(string(message), playerId)

	}
}

func (controller *PlayersController) interpretWebSocketMessage(message string, playerId int) {
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
		case "id":
			controller.sendPlayerId(playerId)
		}
	}
}

func (controller *PlayersController) sendPlayerId(playerId int) {
	controller.connectionsMutex.Lock()
	conn := controller.connections[playerId]
	controller.connectionsMutex.Unlock()
	if conn == nil {
		return
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("id %v", playerId)))
	if err != nil {
		log.Println("Error writing message: ", err)
		return
	}
}

func (controller *PlayersController) getPlayers(c *gin.Context) {
	allPlayers := controller.playersBank.GetAllSearchingPlayers()
	c.JSON(http.StatusOK, allPlayers)
}
