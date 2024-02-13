package Models

import "sync"

// SearchingPlayersList manages a list of searching players.
type SearchingPlayersList struct {
	players      map[int]SearchingPlayer
	nextID       int
	contentMutex sync.Mutex
}

// NewSearchingPlayersList creates a new SearchingPlayersList instance.
func NewSearchingPlayersList() *SearchingPlayersList {
	return &SearchingPlayersList{
		players: map[int]SearchingPlayer{},
		nextID:  1,
	}
}

// createNewPlayerID generates a new player ID.
func (pb *SearchingPlayersList) createNewPlayerID() int {
	newPlayerID := pb.nextID
	pb.nextID++
	return newPlayerID
}

// CreateSearchingPlayer creates a new searching player and returns its ID.
func (pb *SearchingPlayersList) CreateSearchingPlayer() int {
	pb.contentMutex.Lock()
	defer pb.contentMutex.Unlock()
	newPlayerID := pb.createNewPlayerID()
	pb.players[newPlayerID] = NewSearchingPlayer(newPlayerID)
	return newPlayerID
}

// GetAllSearchingPlayers returns all searching players.
func (pb *SearchingPlayersList) GetAllSearchingPlayers() []SearchingPlayer {
	pb.contentMutex.Lock()
	defer pb.contentMutex.Unlock()
	players := make([]SearchingPlayer, 0, len(pb.players))
	for _, v := range pb.players {
		players = append(players, v)
	}
	return players
}

// GetSearchingPlayerFromID returns a searching player with the given ID.
func (pb *SearchingPlayersList) GetSearchingPlayerFromID(id int) *SearchingPlayer {
	pb.contentMutex.Lock()
	defer pb.contentMutex.Unlock()
	player, exists := pb.players[id]
	if exists {
		return &player
	}
	return nil
}

// UpdateSearchingPlayerMetadata updates the metadata of a searching player.
func (pb *SearchingPlayersList) UpdateSearchingPlayerMetadata(id int, key, value string) bool {
	pb.contentMutex.Lock()
	defer pb.contentMutex.Unlock()
	if player, exists := pb.players[id]; exists {
		player.MetaData[key] = value
		return true
	}
	return false
}

// DeleteSearchingPlayer deletes a searching player with the given ID.
func (pb *SearchingPlayersList) DeleteSearchingPlayer(id int) bool {
	pb.contentMutex.Lock()
	defer pb.contentMutex.Unlock()
	_, exists := pb.players[id]
	if exists {
		delete(pb.players, id)
		return true
	}
	return false
}
