package App

import (
	"syndya/internal/AppEnv"
	"syndya/internal/Controllers"
)

// RouteApp create all controller and route for the app
func RouteApp(app *App) {
	Controllers.NewPlayersController(
		app.PlayersBank,
		AppEnv.AppEnv.GetMetadataList(),
		AppEnv.AppEnv.METADATAS_REVIVEPERIOD,
	).Route(app.Router)
}
