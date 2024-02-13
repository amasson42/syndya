package App

import "syndya/internal/Controllers"

// RouteApp create all controller and route for the app
func RouteApp(app *App) {
	Controllers.NewPlayersController(
		app.PlayersBank,
		AppEnv.GetMetadataList(),
		AppEnv.METADATAS_REVIVEPERIOD,
	).Route(app.Router)
}
