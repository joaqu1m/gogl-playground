package main

import (
	"math"
	"runtime"

	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/engine"
	"github.com/joaqu1m/gogl-playground/libs/entities"
	"github.com/joaqu1m/gogl-playground/libs/logger"
	gmath "github.com/joaqu1m/gogl-playground/math"
)

const (
	width  = 800
	height = 600
)

func init() {
	// OpenGL exige que tudo rode na mesma thread do SO
	runtime.LockOSThread()
}

// os calculos de rotacao são feitos em rad, QuatFromAxisAngle faz isso converte vec e  quaterions o float32(math.Pi/2) para quaternions
// é melhor trabalhar com radianos é menos problemas matematicos
func main() {

	app := engine.NewApp(width, height, "OpenGL 4.1 Playground")

	models := []model.Model{
		model.NewModel(
			"Demogorgon",
			"assets/dead_by_daylight_-_eleven.glb",
			entities.Transform{
				Position: gmath.Vec3{X: 0, Y: 0, Z: 0},
				Rotation: gmath.QuatFromAxisAngle(
					gmath.Vec3{X: 0, Y: 0, Z: 0},
					float32(math.Pi/2),
				),
				Scale: gmath.Vec3{X: 0.01, Y: 0.01, Z: 0.01},
			},
		),
		model.NewModel(
			"Shield",
			"assets/shield.glb",
			entities.Transform{
				Position: gmath.Vec3{X: 2, Y: 0, Z: 0},
				Rotation: gmath.QuatFromAxisAngle(
					gmath.Vec3{X: 0, Y: 0, Z: 0},
					float32(math.Pi/2),
				),
				Scale: gmath.Vec3{X: 1, Y: 1, Z: 1},
			},
		),
	}

	app.SetModels(models)
	app.Run()

	logger.Infof("Game closed")
}
