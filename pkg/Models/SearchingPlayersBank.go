package Models

type SearchingPlayersBank interface {
	CreateSearchingPlayer() int
	GetSearchingPlayerFromID(id int) *SearchingPlayer
	GetAllSearchingPlayers() []SearchingPlayer
	UpdateSearchingPlayerMetadata(id int, key string, value string) bool
	DeleteSearchingPlayer(id int) bool
	ForEach(func(*SearchingPlayer))
}
