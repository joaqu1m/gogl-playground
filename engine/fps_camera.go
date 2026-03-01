package engine

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/gmath"
)

// CameraMovement directions for keyboard input.
type CameraMovement int

const (
	Forward CameraMovement = iota
	Backward
	Left
	Right
	Upward
	Downward
)

// FPSCamera implements a first-person camera.  It keeps track of
// yaw/pitch angles and updates position/axes accordingly.  The
// methods satisfy the engine.Camera interface.
type FPSCamera struct {
	Window *glfw.Window

	// Pos holds the world-space position.  named differently to avoid
	// collision with the Position() method required by the Camera
	// interface.
	Pos gmath.Vec3

	Yaw   float32
	Pitch float32

	Front   gmath.Vec3
	Up      gmath.Vec3
	Right   gmath.Vec3
	WorldUp gmath.Vec3

	Fov  float32
	Near float32
	Far  float32

	Speed       float32
	Sensitivity float32

	lastX      float64
	lastY      float64
	firstMouse bool
}

// NewFPSCamera creates a camera positioned at (0,0,3) looking toward
// -Z with a standard up vector.  It also hooks the mouse cursor
// mode and callbacks on the provided window.
func NewFPSCamera(window *glfw.Window) *FPSCamera {
	c := &FPSCamera{
		Window:      window,
		Pos:         gmath.Vec3{X: 0, Y: 0, Z: 3},
		Yaw:         -90, // facing -Z
		Pitch:       0,
		Front:       gmath.Vec3{X: 0, Y: 0, Z: -1},
		WorldUp:     gmath.Vec3{X: 0, Y: 1, Z: 0},
		Fov:         45,
		Near:        0.1,
		Far:         100,
		Speed:       5,
		Sensitivity: 0.1,
		firstMouse:  true,
	}
	c.updateCameraVectors()

	// hide and grab cursor so we can receive deltas
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		c.MouseCallback(xpos, ypos)
	})
	window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		c.ProcessMouseScroll(yoff)
	})

	return c
}

// ViewMatrix returns a look‑at matrix from the camera's current
// position and orientation.
func (c *FPSCamera) ViewMatrix() gmath.Mat4 {
	eye := [3]float32{c.Pos.X, c.Pos.Y, c.Pos.Z}
	center := [3]float32{c.Pos.X + c.Front.X, c.Pos.Y + c.Front.Y, c.Pos.Z + c.Front.Z}
	up := [3]float32{c.Up.X, c.Up.Y, c.Up.Z}
	return gmath.MatLookAt(eye, center, up)
}

// ProjectionMatrix returns a perspective projection based on the
// camera's fov, near and far planes and the supplied aspect ratio.
func (c *FPSCamera) ProjectionMatrix(aspect float32) gmath.Mat4 {
	return gmath.MatPerspective(c.Fov*math.Pi/180, aspect, c.Near, c.Far)
}

// Position returns the current world-space position of the camera.
func (c *FPSCamera) Position() gmath.Vec3 {
	return c.Pos
}

// Update handles keyboard input and modifies the position.
func (c *FPSCamera) Update(dt float32) {
	speed := c.Speed * dt
	if c.Window.GetKey(glfw.KeyW) == glfw.Press {
		c.Pos = c.Pos.Add(c.Front.MulScalar(speed))
	}
	if c.Window.GetKey(glfw.KeyS) == glfw.Press {
		c.Pos = c.Pos.Add(c.Front.MulScalar(-speed))
	}
	if c.Window.GetKey(glfw.KeyA) == glfw.Press {
		c.Pos = c.Pos.Add(c.Right.MulScalar(-speed))
	}
	if c.Window.GetKey(glfw.KeyD) == glfw.Press {
		c.Pos = c.Pos.Add(c.Right.MulScalar(speed))
	}
	if c.Window.GetKey(glfw.KeySpace) == glfw.Press {
		c.Pos = c.Pos.Add(c.WorldUp.MulScalar(speed))
	}
	if c.Window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		c.Pos = c.Pos.Add(c.WorldUp.MulScalar(-speed))
	}
}

// MouseCallback should be called whenever the cursor moves.  It
// updates the yaw/pitch angles and recomputes the camera axes.
func (c *FPSCamera) MouseCallback(xpos, ypos float64) {
	if c.firstMouse {
		c.lastX = xpos
		c.lastY = ypos
		c.firstMouse = false
	}

	xoffset := float32(xpos - c.lastX)
	yoffset := float32(c.lastY - ypos) // reversed since y-coordinates go from bottom to top
	c.lastX = xpos
	c.lastY = ypos

	xoffset *= c.Sensitivity
	yoffset *= c.Sensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	if c.Pitch > 89 {
		c.Pitch = 89
	}
	if c.Pitch < -89 {
		c.Pitch = -89
	}

	c.updateCameraVectors()
}

// ProcessMouseScroll changes the field of view (zoom).
func (c *FPSCamera) ProcessMouseScroll(yoffset float64) {
	c.Fov -= float32(yoffset)
	if c.Fov < 1 {
		c.Fov = 1
	}
	if c.Fov > 45 {
		c.Fov = 45
	}
}

// updateCameraVectors recalculates the front, right and up vectors
// from the current yaw and pitch values.
func (c *FPSCamera) updateCameraVectors() {
	yawRad := c.Yaw * math.Pi / 180
	pitchRad := c.Pitch * math.Pi / 180
	front := gmath.Vec3{
		X: float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad))),
		Y: float32(math.Sin(float64(pitchRad))),
		Z: float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad))),
	}
	c.Front = front.Normalize()
	c.Right = gmath.Cross(c.Front, c.WorldUp).Normalize()
	c.Up = gmath.Cross(c.Right, c.Front).Normalize()
}
