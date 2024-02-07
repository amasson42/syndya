package syndya

import (
	"syndya/internal/App"
	"syndya/internal/Controllers"
)

func Run(listener string) {
	app := App.MakeApp()
	Controllers.RouteApp(app)
	app.Router.Run(listener)
	TearDown(app)
}

func TearDown(app *App.App) {

}
