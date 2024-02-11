package Models

import "time"

type SearchingPlayer struct {
	ID        int            `json:"id"`
	TimeStamp int64          `json:"timestamp"`
	MetaData  map[string]any `json:"metadata"`
}

func NewSearchingPlayer(id int) SearchingPlayer {
	return SearchingPlayer{
		ID:        id,
		TimeStamp: time.Now().Unix(),
		MetaData:  map[string]any{},
	}
}
