package Controllers

import (
	"syndya/internal/App"
)

func RouteApp(app *App.App) {
	NewPlayersController().Route(app.Router)
}
