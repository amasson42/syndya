package syndya

import (
	"syndya/internal/App"
)

func Run(listener string) {
	app := App.MakeApp()
	App.RouteApp(app)
	app.Router.Run(listener)
	TearDown(app)
}

func TearDown(app *App.App) {

}
