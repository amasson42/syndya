package App

import (
	"log"
	"syndya/internal/AppEnv"
	"syndya/internal/MatchFinder"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router      *gin.Engine
	PlayersBank Models.SearchingPlayersBank
	MatchFinder *MatchFinder.MatchFinder
}

func MakeApp() *App {
	playersList := Models.NewSearchingPlayersList()
	app := App{
		Router:      gin.Default(),
		PlayersBank: playersList,
	}
	return &app
}

func (app *App) MatchFinderService() {
	if !AppEnv.AppEnv.HasMatchFinderScript() {
		return
	}
	matchfinder, err := MatchFinder.NewMatchFinder(
		app.PlayersBank,
		AppEnv.AppEnv.MATCHFINDER_LUASCRIPT,
		AppEnv.AppEnv.MATCHFINDER_RESETSTATE,
	)
	if err != nil {
		log.Println("Error loading matchfinder script: ", err)
		return
	}

	app.MatchFinder = matchfinder

	matchfinder.MatchupDelegate = func(ids []int) {
		log.Printf("We have a MATCH ! %v\n", ids)
	}

	matchfinder.AsyncRunLoop(AppEnv.AppEnv.MATCHFINDER_TIMEINTERVAL)
}
