package Controllers

import (
	"fmt"
	"log"
	"strings"
	"syndya/internal/App"
	"syndya/pkg/Models"

	"github.com/gorilla/websocket"
)

type PlayerConnection struct {
	playerId        int
	conn            *websocket.Conn
	playersBank     Models.SearchingPlayersBank
	cachedMetaDatas map[string]string
}

func NewPlayerConnection(playerId int, conn *websocket.Conn, playersBank Models.SearchingPlayersBank) *PlayerConnection {
	return &PlayerConnection{
		playerId:        playerId,
		conn:            conn,
		playersBank:     playersBank,
		cachedMetaDatas: map[string]string{},
	}
}

func (playerConnection *PlayerConnection) SendPlayerId() {
	if playerConnection.conn == nil {
		return
	}
	err := playerConnection.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("id %v", playerConnection.playerId)))
	if err != nil {
		log.Println("Error writing message: ", err)
		return
	}
}

func (playerConnection *PlayerConnection) InterpretWebSocketMessage(message string) {
	cmd := strings.SplitN(string(message[:]), " ", 2)

	if len(cmd) >= 1 {
		switch cmd[0] {
		case "meta":
			if len(cmd) >= 2 {
				parameters := strings.SplitN(cmd[1], "=", 2)
				if len(parameters) == 2 {
					metaKey := parameters[0]
					metaValue := parameters[1]
					playerConnection.playersBank.UpdateSearchingPlayerMetadata(playerConnection.playerId, metaKey, metaValue)
					playerConnection.cachedMetaDatas[metaKey] = metaValue
				}
			}
		case "id":
			playerConnection.SendPlayerId()
		}
	}
}

func (playerConnection *PlayerConnection) RequestMissingMetadatas() bool {
	if App.AppEnv.METADATAS_LIST == "_" {
		return true
	}
	list := strings.Split(App.AppEnv.METADATAS_LIST, ",")
	requestedNone := true
	for _, metaName := range list {
		if _, isMapContainsKey := playerConnection.cachedMetaDatas[metaName]; !isMapContainsKey {
			playerConnection.RequestMetadata(metaName)
			requestedNone = false
		}
	}
	return requestedNone
}

func (playerConnection *PlayerConnection) RequestMetadata(key string) {
	err := playerConnection.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("meta %v", key)))
	if err != nil {
		log.Println("Error writing message: ", err)
		return
	}
}
