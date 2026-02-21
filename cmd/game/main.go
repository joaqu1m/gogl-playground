package main

import (
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/libs/entities"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

const (
	width  = 800
	height = 600
)

// ModelEntry define um modelo a ser carregado na cena.
type ModelEntry struct {
	Path        string
	Scale       float32
	Translation [3]float32
	loaded      *gltfloader.GLTFModel
}

var (
	shaderProgram uint32
	timeAccum     float64

	models = []model.Model{}
)

func init() {
	// OpenGL exige que tudo rode na mesma thread do SO
	runtime.LockOSThread()
}

func main() {
	initGLFW()
	defer glfw.Terminate()

	window := createWindow()
	initOpenGL()

	models = []model.Model{
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

	shaderProgram = createShaderProgram()

	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		timeAccum = currentTime - previousTime
		_ = previousTime
		previousTime = currentTime

		draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	logger.Infof("Exiting game loop")
}

// ---------------- GLFW ----------------

func initGLFW() {
	if err := glfw.Init(); err != nil {
		logger.Fatalf("failed to init glfw: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func createWindow() *glfw.Window {
	window, err := glfw.CreateWindow(width, height, "OpenGL 4.1 Playground", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

// ---------------- OpenGL ----------------

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	logger.Infof("OpenGL version: %s", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
}

// ---------------- Shaders ----------------

var vertexShaderSource = `#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 vNormal;
out vec3 vFragPos;
out vec2 vTexCoord;

void main() {
	vFragPos = vec3(model * vec4(aPos, 1.0));
	vNormal = mat3(transpose(inverse(model))) * aNormal;
	vTexCoord = aTexCoord;
	gl_Position = projection * view * vec4(vFragPos, 1.0);
}` + "\x00"

var fragmentShaderSource = `#version 410 core
in vec3 vNormal;
in vec3 vFragPos;
in vec2 vTexCoord;

out vec4 FragColor;

uniform vec3 lightDir;
uniform sampler2D diffuseMap;
uniform vec4 baseColor;
uniform int useTexture;

void main() {
	// Cor base: textura ou cor do material
	vec3 color;
	if (useTexture == 1) {
		color = texture(diffuseMap, vTexCoord).rgb * baseColor.rgb;
	} else {
		color = baseColor.rgb;
	}

	// Ambient
	float ambientStrength = 0.2;
	vec3 ambient = ambientStrength * vec3(1.0);

	// Diffuse
	vec3 norm = normalize(vNormal);
	vec3 light = normalize(-lightDir);
	float diff = max(dot(norm, light), 0.0);
	vec3 diffuse = diff * vec3(1.0);

	vec3 result = (ambient + diffuse) * color;
	FragColor = vec4(result, baseColor.a);
}` + "\x00"

func createShaderProgram() uint32 {
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		panic(string(log))
	}

	return shader
}

// ---------------- Render ----------------

// angle acumula a rotação ao longo do tempo
var angle float64

func draw() {
	angle += 1.0 * timeAccum

	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shaderProgram)

	// Rotação global compartilhada por todos os modelos
	rotMat := matRotateY(float32(angle))

	viewMat := matLookAt(
		[3]float32{0, 0.8, 3.0},
		[3]float32{0, 0, 0},
		[3]float32{0, 1, 0},
	)
	projMat := matPerspective(
		float32(45.0*math.Pi/180.0),
		float32(width)/float32(height),
		0.1, 100.0,
	)

	setUniformMat4(shaderProgram, "view", viewMat)
	setUniformMat4(shaderProgram, "projection", projMat)
	setUniformVec3(shaderProgram, "lightDir", [3]float32{-0.3, -0.8, -0.5})

	// Renderiza cada modelo da cena
	for _, entry := range models {
		// Transform por modelo: translate * rotate * scale
		scaleMat := matScale(entry.Scale.Width, entry.Scale.Height, entry.Scale.Depth)
		transMat := matTranslate(entry.Translation.ToArray())
		baseMat := matMul(transMat, matMul(rotMat, scaleMat))

		logger.Debugf(entry.Name)
		for _, m := range entry.LoadedModel.Meshes {
			modelMat := matMul(baseMat, mat4(m.Transform))
			setUniformMat4(shaderProgram, "model", modelMat)

			// Material
			bcLoc := gl.GetUniformLocation(shaderProgram, gl.Str("baseColor\x00"))
			gl.Uniform4f(bcLoc, m.BaseColor[0], m.BaseColor[1], m.BaseColor[2], m.BaseColor[3])

			utLoc := gl.GetUniformLocation(shaderProgram, gl.Str("useTexture\x00"))
			if m.HasTexture {
				gl.Uniform1i(utLoc, 1)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, m.TextureID)
				dmLoc := gl.GetUniformLocation(shaderProgram, gl.Str("diffuseMap\x00"))
				gl.Uniform1i(dmLoc, 0)
			} else {
				gl.Uniform1i(utLoc, 0)
			}

			gl.BindVertexArray(m.VAO)
			if m.HasIndices {
				gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
			} else {
				gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
			}
		}
	}
}

// ---------------- Math (matrizes column-major para OpenGL) ----------------

// mat4 é column-major [16]float32, como OpenGL espera.
type mat4 [16]float32

func matIdentity() mat4 {
	return mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func matRotateY(angle float32) mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	return mat4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

func matMul(a, b mat4) mat4 {
	var r mat4
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float32
			for k := 0; k < 4; k++ {
				sum += a[k*4+row] * b[col*4+k]
			}
			r[col*4+row] = sum
		}
	}
	return r
}

func matScaleUniform(s float32) mat4 {
	return matScale(s, s, s)
}

func matScale(sx, sy, sz float32) mat4 {
	return mat4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

func matTranslate(t [3]float32) mat4 {
	return mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		t[0], t[1], t[2], 1,
	}
}

func matPerspective(fovy, aspect, near, far float32) mat4 {
	f := float32(1.0 / math.Tan(float64(fovy/2.0)))
	nf := near - far
	return mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / nf, -1,
		0, 0, (2 * far * near) / nf, 0,
	}
}

func matLookAt(eye, center, up [3]float32) mat4 {
	f := vecNormalize(vecSub(center, eye))
	s := vecNormalize(vecCross(f, up))
	u := vecCross(s, f)

	return mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		-vecDot(s, eye), -vecDot(u, eye), vecDot(f, eye), 1,
	}
}

func vecSub(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func vecCross(a, b [3]float32) [3]float32 {
	return [3]float32{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func vecDot(a, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func vecNormalize(v [3]float32) [3]float32 {
	l := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
	if l == 0 {
		return v
	}
	return [3]float32{v[0] / l, v[1] / l, v[2] / l}
}

// ---------------- Uniform helpers ----------------

func setUniformMat4(program uint32, name string, m mat4) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &m[0])
}

func setUniformVec3(program uint32, name string, v [3]float32) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform3f(loc, v[0], v[1], v[2])
}
