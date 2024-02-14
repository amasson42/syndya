package Models

import "time"

type SearchingPlayer struct {
	ID        int               `json:"id"`
	TimeStamp int64             `json:"timestamp"`
	MetaData  map[string]string `json:"metadata"`
	Complete  bool              `json:"complete"`
	GameAddr  *string           `json:"gameaddr"`
}

func NewSearchingPlayer(id int) SearchingPlayer {
	return SearchingPlayer{
		ID:        id,
		TimeStamp: time.Now().Unix(),
		MetaData:  map[string]string{},
		Complete:  false,
		GameAddr:  nil,
	}
}
