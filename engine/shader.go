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

uniform sampler2D diffuseMap;
uniform vec4 baseColor;
uniform int useTexture;

uniform vec3 viewPos;
uniform float specularStrength;
uniform float shininess;

#define MAX_LIGHTS 8
#define LIGHT_DIRECTIONAL 0
#define LIGHT_POINT       1
#define LIGHT_SPOT        2

struct Light {
	int type;
	vec3 position;
	vec3 direction;
	vec3 color;
	float ambient;
	float constant;
	float linear;
	float quadratic;
	float cutOff;
	float outerCutOff;
};

uniform Light lights[MAX_LIGHTS];
uniform int numLights;

vec3 calcDirectional(Light l, vec3 norm, vec3 viewDir) {
	vec3 lightDir = normalize(-l.direction);

	vec3 ambient = l.ambient * l.color;

	float diff   = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * l.color;

	vec3 halfDir  = normalize(lightDir + viewDir);
	float spec    = pow(max(dot(norm, halfDir), 0.0), shininess);
	vec3 specular = specularStrength * spec * l.color;

	return ambient + diffuse + specular;
}

vec3 calcPoint(Light l, vec3 norm, vec3 viewDir) {
	vec3 lightDir = normalize(l.position - vFragPos);

	float dist        = length(l.position - vFragPos);
	float attenuation = 1.0 / (l.constant + l.linear * dist + l.quadratic * dist * dist);

	vec3 ambient = l.ambient * l.color;

	float diff   = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * l.color;

	vec3 halfDir  = normalize(lightDir + viewDir);
	float spec    = pow(max(dot(norm, halfDir), 0.0), shininess);
	vec3 specular = specularStrength * spec * l.color;

	return (ambient + diffuse + specular) * attenuation;
}

vec3 calcSpot(Light l, vec3 norm, vec3 viewDir) {
	vec3 lightDir = normalize(l.position - vFragPos);

	float theta     = dot(lightDir, normalize(-l.direction));
	float epsilon   = l.cutOff - l.outerCutOff;
	float intensity = clamp((theta - l.outerCutOff) / epsilon, 0.0, 1.0);

	float dist        = length(l.position - vFragPos);
	float attenuation = 1.0 / (l.constant + l.linear * dist + l.quadratic * dist * dist);

	vec3 ambient = l.ambient * l.color;

	float diff   = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * l.color;

	vec3 halfDir  = normalize(lightDir + viewDir);
	float spec    = pow(max(dot(norm, halfDir), 0.0), shininess);
	vec3 specular = specularStrength * spec * l.color;

	diffuse  *= intensity;
	specular *= intensity;

	return (ambient + diffuse + specular) * attenuation;
}

void main() {
	vec3 color;
	if (useTexture == 1) {
		color = texture(diffuseMap, vTexCoord).rgb * baseColor.rgb;
	} else {
		color = baseColor.rgb;
	}

	vec3 norm    = normalize(vNormal);
	vec3 viewDir = normalize(viewPos - vFragPos);

	vec3 result = vec3(0.0);
	for (int i = 0; i < numLights; i++) {
		if (lights[i].type == LIGHT_DIRECTIONAL) {
			result += calcDirectional(lights[i], norm, viewDir);
		} else if (lights[i].type == LIGHT_POINT) {
			result += calcPoint(lights[i], norm, viewDir);
		} else if (lights[i].type == LIGHT_SPOT) {
			result += calcSpot(lights[i], norm, viewDir);
		}
	}

	FragColor = vec4(result * color, baseColor.a);
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
