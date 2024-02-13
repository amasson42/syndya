package Controllers

import (
	"log"
	"net/http"
	"sync"
	"syndya/internal/App"
	"syndya/pkg/Models"
	"time"

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

	playerConnection := NewPlayerConnection(conn, controller.playersBank)

	controller.connectionsMutex.Lock()
	controller.connections[playerConnection.playerId] = playerConnection
	controller.connectionsMutex.Unlock()

	terminateParrallelRoutine := make(chan struct{})

	defer func() {
		close(terminateParrallelRoutine)
		controller.connectionsMutex.Lock()
		delete(controller.connections, playerConnection.playerId)
		controller.connectionsMutex.Unlock()
		playerConnection.Close()
	}()

	playerConnection.SendPlayerId()

	playerConnection.RequestMissingMetadatas()

	requestMetadatasTicker := time.NewTicker(time.Duration(App.AppEnv.METADATAS_REVIVEPERIOD) * time.Millisecond)

	go func() {
		for {
			select {
			case <-requestMetadatasTicker.C:
				finished := playerConnection.RequestMissingMetadatas()
				if finished {
					requestMetadatasTicker.Stop()
				}
			case <-terminateParrallelRoutine:
				return
			}
		}
	}()

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
