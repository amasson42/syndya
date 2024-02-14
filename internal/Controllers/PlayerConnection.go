package Controllers

import (
	"fmt"
	"log"
	"strings"
	"syndya/pkg/Models"

	"github.com/gorilla/websocket"
)

// PlayerConnection represents a connection with a player.
type PlayerConnection struct {
	playerId        int
	conn            *websocket.Conn
	playersBank     Models.SearchingPlayersBank
	cachedMetaDatas map[string]string
}

// NewPlayerConnection creates a new PlayerConnection instance.
func NewPlayerConnection(conn *websocket.Conn, playersBank Models.SearchingPlayersBank) *PlayerConnection {
	playerId := playersBank.CreateSearchingPlayer()

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Closed connection with %v\n", playerId)
		return nil
	})

	return &PlayerConnection{
		playerId:        playerId,
		conn:            conn,
		playersBank:     playersBank,
		cachedMetaDatas: make(map[string]string),
	}
}

// Close closes the connection and removes the player from the players bank.
func (pc *PlayerConnection) Close() {
	pc.conn.Close()
	pc.playersBank.DeleteSearchingPlayer(pc.playerId)
}

// SendPlayerId sends the player's ID over the connection.
func (pc *PlayerConnection) SendPlayerId() {
	pc.sendMessage(fmt.Sprintf("id %v", pc.playerId))
}

// InterpretWebSocketMessage interprets a message received over the WebSocket connection.
func (pc *PlayerConnection) InterpretWebSocketMessage(message string) {
	cmd := strings.SplitN(string(message[:]), " ", 2)

	if len(cmd) < 1 {
		return
	}

	switch cmd[0] {
	case "meta":
		pc.handleMetadataCommand(cmd)
	case "id":
		pc.SendPlayerId()
	}
}

// handleMetadataCommand handles the "meta" command received over WebSocket.
func (pc *PlayerConnection) handleMetadataCommand(cmd []string) {
	if len(cmd) < 2 {
		return
	}

	parameters := strings.SplitN(cmd[1], "=", 2)
	if len(parameters) != 2 {
		return
	}

	pc.updateMetadata(parameters[0], parameters[1])
}

// updateMetadata updates the metadata for the player.
func (pc *PlayerConnection) updateMetadata(key string, value string) {
	pc.playersBank.UpdateSearchingPlayerMetadata(pc.playerId, key, value)
	pc.cachedMetaDatas[key] = value
}

// RequestMissingMetadatas requests missing metadata from the player.
func (pc *PlayerConnection) RequestMissingMetadatas(datalist []string) bool {
	requestedNone := true
	for _, metaName := range datalist {
		if _, keyExists := pc.cachedMetaDatas[metaName]; !keyExists {
			pc.RequestMetadata(metaName)
			requestedNone = false
		}
	}
	return requestedNone
}

// RequestMetadata requests a specific metadata from the player.
func (pc *PlayerConnection) RequestMetadata(key string) {
	pc.sendMessage(fmt.Sprintf("meta %v", key))
}

// SendGameAddr send the game address to the player.
func (pc *PlayerConnection) SendGameAddr(addr string) {
	pc.sendMessage(fmt.Sprintf("game %v", addr))
}

// sendMessage sends a message over the WebSocket connection.
func (pc *PlayerConnection) sendMessage(message string) {
	if pc.conn == nil {
		return
	}
	err := pc.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error writing message: ", err)
		return
	}
}
