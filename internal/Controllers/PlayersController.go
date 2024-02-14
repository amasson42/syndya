package Controllers

import (
	"log"
	"net/http"
	"sync"
	"syndya/pkg/Models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// PlayersController handles WebSocket connections from players.
type PlayersController struct {
	upgrader             websocket.Upgrader
	playersBank          Models.SearchingPlayersBank
	connections          map[int]*PlayerConnection
	connectionsMutex     sync.Mutex
	metadataList         []string
	metadataRevivePeriod int
}

// NewPlayersController creates a new instance of PlayersController.
func NewPlayersController(
	playersBank Models.SearchingPlayersBank,
	metadataList []string,
	metadataRevivePeriod int,
) *PlayersController {
	return &PlayersController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		playersBank:          playersBank,
		connections:          make(map[int]*PlayerConnection),
		metadataList:         metadataList,
		metadataRevivePeriod: metadataRevivePeriod,
	}
}

// Route routes the endpoints for PlayersController.
func (controller *PlayersController) Route(router *gin.Engine) {
	router.GET("/search", controller.searchGame)
	router.GET("/players", controller.getPlayers)
}

// searchGame handles WebSocket connections for searching games.
func (controller *PlayersController) searchGame(c *gin.Context) {
	conn, err := controller.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading request: ", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	playerConnection := NewPlayerConnection(conn, controller.playersBank)
	defer playerConnection.Close()

	controller.addConnection(playerConnection)
	defer controller.removeConnection(playerConnection)

	playerConnection.SendPlayerId()
	playerConnection.RequestMissingMetadatas(controller.metadataList)

	terminateParrallelRoutine := make(chan struct{})
	defer close(terminateParrallelRoutine)

	go controller.monitorMetadataRequests(playerConnection, terminateParrallelRoutine)

	controller.readMessages(playerConnection)
}

// addConnection adds a player connection to the controller's connections map.
func (controller *PlayersController) addConnection(pc *PlayerConnection) {
	controller.connectionsMutex.Lock()
	defer controller.connectionsMutex.Unlock()
	controller.connections[pc.playerId] = pc
}

// removeConnection removes a player connection from the controller's connections map.
func (controller *PlayersController) removeConnection(pc *PlayerConnection) {
	controller.connectionsMutex.Lock()
	defer controller.connectionsMutex.Unlock()
	delete(controller.connections, pc.playerId)
}

// monitorMetadataRequests continuously requests missing metadata for the player.
func (controller *PlayersController) monitorMetadataRequests(pc *PlayerConnection, terminate <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(controller.metadataRevivePeriod) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			finished := pc.RequestMissingMetadatas(controller.metadataList)
			if finished {
				controller.playersBank.SetSearchingPlayerComplete(pc.playerId, true)
				return
			}
		case <-terminate:
			return
		}
	}
}

// readMessages reads messages from the player connection and interprets them.
func (controller *PlayersController) readMessages(pc *PlayerConnection) {
	for {
		_, message, err := pc.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		pc.InterpretWebSocketMessage(string(message))
	}
}

func (controller *PlayersController) getPlayers(c *gin.Context) {
	allPlayers := controller.playersBank.GetAllSearchingPlayers()
	c.JSON(http.StatusOK, allPlayers)
}
