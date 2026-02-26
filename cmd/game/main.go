package main

import (
	"runtime"

	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/engine"
	"github.com/joaqu1m/gogl-playground/libs/entities"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

const (
	width  = 800
	height = 600
)

func init() {
	// OpenGL exige que tudo rode na mesma thread do SO
	runtime.LockOSThread()
}

func main() {
	app := engine.NewApp(width, height, "OpenGL 4.1 Playground")

	models := []model.Model{
		model.NewModel(
			"Demogorgon",
			"assets/dead_by_daylight_-_eleven.glb",
			entities.Dimensions{Width: 0.01, Height: 0.01, Depth: 0.01},
			entities.Dimensions{Width: 0, Height: 0, Depth: 0},
		),
		model.NewModel(
			"Shield",
			"assets/shield.glb",
			entities.Dimensions{Width: 1.0, Height: 1.0, Depth: 1.0},
			entities.Dimensions{Width: 2.0, Height: 0, Depth: 0},
		),
	}

	app.SetModels(models)
	app.Run()

	logger.Infof("Game closed")
}
