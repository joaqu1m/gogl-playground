package main

import (
	"math"
	"runtime"

	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/engine"
	"github.com/joaqu1m/gogl-playground/gmath"
	"github.com/joaqu1m/gogl-playground/libs/entities"
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

	// create application and automatically set up an FPS camera
	// The camera will receive mouse movement events (cursor is
	// disabled/hidden) and handle WASD+space/shift for translation.
	app := engine.NewApp(width, height, "OpenGL 4.1 Playground")

	// if you ever need to access the concrete FPSCamera instance
	// you can assert the interface and call its methods directly:
	//
	//     cam := app.Camera.(*engine.FPSCamera)
	//     cam.Pos = gmath.Vec3{X:0,Y:1,Z:5} // field is now "Pos"
	//     cam.Yaw = 45
	//     cam.updateCameraVectors()
	//
	// the mouse callback that updates yaw/pitch is already
	// registered by NewFPSCamera above, so moving the real mouse
	// in the window will make the scene look around.

	app.ModelManager.AddModel(model.NewModel(
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
	))

	app.ModelManager.AddModel(model.NewModel(
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
	))

	app.Run()
}
