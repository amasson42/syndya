package Models

type SearchingPlayersList struct {
	players map[int]SearchingPlayer
	nextId  int
}

func NewSearchingPlayersList() *SearchingPlayersList {
	return &SearchingPlayersList{
		players: map[int]SearchingPlayer{},
		nextId:  1,
	}
}

func (playersBank *SearchingPlayersList) createNewPlayerId() int {
	newPlayerId := playersBank.nextId
	playersBank.nextId += 1
	return newPlayerId
}

func (playersBank *SearchingPlayersList) CreateSearchingPlayer() int {
	newPlayerId := playersBank.createNewPlayerId()
	playersBank.players[newPlayerId] = NewSearchingPlayer(newPlayerId)
	return newPlayerId
}

func (playersBank *SearchingPlayersList) GetAllSearchingPlayers() []SearchingPlayer {
	players := []SearchingPlayer{}
	for _, v := range playersBank.players {
		players = append(players, v)
	}
	return players
}

func (playersBank *SearchingPlayersList) GetSearchingPlayerFromID(id int) *SearchingPlayer {
	if player, exists := playersBank.players[id]; exists {
		return &player
	}
	return nil
}

func (playersBank *SearchingPlayersList) UpdateSearchingPlayerMetadata(id int, key string, value string) bool {
	if _, exists := playersBank.players[id]; exists {
		playersBank.players[id].MetaData[key] = value
		return true
	} else {
		return false
	}
}

func (playersBank *SearchingPlayersList) DeleteSearchingPlayer(id int) bool {
	delete(playersBank.players, id)
	return true
}
