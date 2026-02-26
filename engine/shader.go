package engine

import "github.com/go-gl/gl/v4.1-core/gl"

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
