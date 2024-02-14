package syndya

import (
	"syndya/internal/App"
	"syndya/internal/AppEnv"
)

func Run() {
	app := App.MakeApp()
	App.RouteApp(app)

	app.MatchFinderService()

	app.Router.Run(AppEnv.AppEnv.GetListener())

	TearDown(app)
}

func TearDown(app *App.App) {

}
