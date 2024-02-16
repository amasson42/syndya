package App

import (
	"log"
	"syndya/internal/AppEnv"
	"syndya/internal/GameDeployer"
	"syndya/internal/MatchFinder"
	"syndya/pkg/Models"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router       *gin.Engine
	PlayersBank  Models.SearchingPlayersBank
	MatchFinder  *MatchFinder.MatchFinder
	GameDeployer *GameDeployer.GameDeployer
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
		log.Println("No script defined in MATCHFINDER_LUASCRIPT... Matchmaking service will not run")
		return
	}
	matchfinder, err := MatchFinder.NewMatchFinder(
		app.PlayersBank,
		AppEnv.AppEnv.METADATAS_REQUIRED,
		AppEnv.AppEnv.MATCHFINDER_LUASCRIPT,
		AppEnv.AppEnv.MATCHFINDER_RESETSTATE,
	)
	if err != nil {
		log.Println("Error loading matchfinder script: ", err)
		return
	}

	app.MatchFinder = matchfinder

	matchfinder.MatchupDelegate = app.startGameWithPlayers

	matchfinder.AsyncRunLoop(AppEnv.AppEnv.MATCHFINDER_TIMEINTERVAL)
}

func (app *App) startGameWithPlayers(ids []int) {
	if app.GameDeployer != nil {
		players := app.PlayersBank.GetSearchingPlayerFromIDs(ids)
		go func() {
			for _, id := range ids {
				app.PlayersBank.SetSearchingPlayerIsJoiningGame(id, true)
			}
			gameaddr, err := app.GameDeployer.DeployGame(players)
			if err != nil {
				log.Printf("Error deploying a game: %v", err)
				return
			}
			for _, id := range ids {
				app.PlayersBank.SetSearchingPlayerGameAddr(id, *gameaddr)
				app.PlayersBank.SetSearchingPlayerIsJoiningGame(id, false)
			}
		}()
	}
}

func (app *App) GameDeployerService() {
	if !AppEnv.AppEnv.HasGameDeployerScript() {
		log.Println("No script defined in GAMEDEPLOYER_LUASCRIPT... Game deloyement will not run")
		return
	}

	gamedeployer := GameDeployer.NewGameDeployer(
		AppEnv.AppEnv.GAMEDEPLOYER_LUASCRIPT,
	)

	app.GameDeployer = gamedeployer
}
