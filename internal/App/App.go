package App

import (
	"syndya/internal/AppEnv"
	"syndya/internal/MatchFinder"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router      *gin.Engine
	PlayersBank Models.SearchingPlayersBank
}

func MakeApp() *App {
	playersList := Models.NewSearchingPlayersList()
	app := App{
		Router:      gin.Default(),
		PlayersBank: playersList,
	}
	return &app
}

func (app *App) StartMatchFinder() {
	if !AppEnv.AppEnv.HasMatchFinderScript() {
		return
	}
	matchfinder := MatchFinder.NewMatchFinder(
		app.PlayersBank,
		AppEnv.AppEnv.MATCHFINDER_LUASCRIPT,
	)
	matchfinder.AsyncRunLoop(AppEnv.AppEnv.MATCHFINDER_TIMEINTERVAL)
}
