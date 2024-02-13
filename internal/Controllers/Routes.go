package Controllers

import (
	"syndya/internal/App"
)

// RouteApp create all controller and route for the app
func RouteApp(app *App.App) {
	NewPlayersController(app.PlayersBank).Route(app.Router)
}
