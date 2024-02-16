package Models

type SearchingPlayersBank interface {
	CreateSearchingPlayer() int
	GetSearchingPlayerFromID(id int) *SearchingPlayer
	GetSearchingPlayerFromIDs(ids []int) []*SearchingPlayer
	GetAllSearchingPlayers() []SearchingPlayer
	UpdateSearchingPlayerMetadata(id int, key string, value string) bool
	SetSearchingPlayerComplete(id int, complete bool) bool
	SetSearchingPlayerIsJoiningGame(id int, joining bool) bool
	SetSearchingPlayerGameAddr(id int, addr string) bool
	DeleteSearchingPlayer(id int) bool
	ForEach(f func(*SearchingPlayer), ignoreIncompletes bool)
}
