package Controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PlayersController struct {
	upgrader         websocket.Upgrader
	playersBank      Models.SearchingPlayersBank
	connections      map[int]*PlayerConnection
	connectionsMutex sync.Mutex
}

func NewPlayersController(playersBank Models.SearchingPlayersBank) *PlayersController {
	controller := PlayersController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		playersBank: playersBank,
		connections: map[int]*PlayerConnection{},
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
	playerConnection := NewPlayerConnection(playerId, conn, controller.playersBank)

	controller.connectionsMutex.Lock()
	controller.connections[playerId] = playerConnection
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

	playerConnection.SendPlayerId()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		playerConnection.InterpretWebSocketMessage(string(message))

	}
}

func (controller *PlayersController) getPlayers(c *gin.Context) {
	allPlayers := controller.playersBank.GetAllSearchingPlayers()
	c.JSON(http.StatusOK, allPlayers)
}
