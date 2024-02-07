package main

import (
	"syndya/internal/App"
	"syndya/internal/syndya"
)

func main() {
	syndya.Run(App.AppEnv.GetListener())
}
