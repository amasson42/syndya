package App

import "github.com/gin-gonic/gin"

type App struct {
	Router *gin.Engine
}

func MakeApp() *App {
	app := App{
		Router: gin.Default(),
	}
	return &app
}
