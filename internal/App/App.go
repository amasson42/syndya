package App

import (
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
