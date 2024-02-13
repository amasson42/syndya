package MatchFinder

import (
	"log"
	"syndya/pkg/Models"
	"time"
)

type MatchFinder struct {
	playersBank Models.SearchingPlayersBank
}

func NewMatchFinder(
	playersBank Models.SearchingPlayersBank,
	scriptPath string,
) *MatchFinder {
	return &MatchFinder{
		playersBank: playersBank,
	}
}

func (mf *MatchFinder) AsyncRunLoop(timeInterval int) {
	go func() {
		for {
			time.Sleep(time.Duration(timeInterval) * time.Millisecond)
			mf.RunOnce()
		}
	}()
}

func (mf *MatchFinder) RunOnce() {
	log.Print("Running match finder loop")
}
