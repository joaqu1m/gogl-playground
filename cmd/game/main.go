package main

import (
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/engine"
	"github.com/joaqu1m/gogl-playground/gmath"
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

// os calculos de rotacao são feitos em rad, QuatFromAxisAngle faz isso converte vec e  quaterions o float32(math.Pi/2) para quaternions
// é melhor trabalhar com radianos é menos problemas matematicos
func main() {

	app := engine.NewApp(width, height, "OpenGL 4.1 Playground")

	app.Models = []model.Model{
		model.NewModel(
			"Eleven",
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

	previousTime := glfw.GetTime()

	for !app.Window.ShouldClose() {

		currentTime := glfw.GetTime()
		app.TimeAccum = currentTime - previousTime
		previousTime = currentTime

		app.Draw()

		app.Window.SwapBuffers()
		glfw.PollEvents()
	}

	logger.Infof("Exiting game loop")
	glfw.Terminate()

	logger.Infof("Game closed")
}
