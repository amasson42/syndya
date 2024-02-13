package Controllers

import (
	"fmt"
	"log"
	"strings"
	"syndya/pkg/Models"

	"github.com/gorilla/websocket"
)

type PlayerConnection struct {
	playerId    int
	conn        *websocket.Conn
	playersBank Models.SearchingPlayersBank
}

func NewPlayerConnection(playerId int, conn *websocket.Conn, playersBank Models.SearchingPlayersBank) *PlayerConnection {
	return &PlayerConnection{
		playerId:    playerId,
		conn:        conn,
		playersBank: playersBank,
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
				}
			}
		case "id":
			playerConnection.SendPlayerId()
		}
	}
}
