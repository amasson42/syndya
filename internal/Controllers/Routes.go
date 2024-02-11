package Controllers

import (
	"syndya/internal/App"
)

func RouteApp(app *App.App) {
	NewPlayersController(app.PlayersBank).Route(app.Router)
}
