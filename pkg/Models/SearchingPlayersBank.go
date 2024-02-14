package Models

type SearchingPlayersBank interface {
	CreateSearchingPlayer() int
	GetSearchingPlayerFromID(id int) *SearchingPlayer
	GetAllSearchingPlayers() []SearchingPlayer
	UpdateSearchingPlayerMetadata(id int, key string, value string) bool
	SetSearchingPlayerComplete(id int, complete bool) bool
	SetSearchingPlayerGameAddr(id int, addr string) bool
	DeleteSearchingPlayer(id int) bool
	ForEach(f func(*SearchingPlayer), ignoreIncompletes bool)
}
