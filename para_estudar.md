# Visão Geral do Projeto

Este é um **playground de engine 3D em Go** usando OpenGL 4.1.

## Arquitetura Geral

```
┌─────────────────────────────────────────────────────────────────┐
│                        main.go (entrada)                        │
├─────────────────────────────────────────────────────────────────┤
│  GLFW (janela/input) → OpenGL (renderização) → Shaders (GPU)   │
├─────────────────────────────────────────────────────────────────┤
│                      gltfloader (modelos 3D)                    │
└─────────────────────────────────────────────────────────────────┘
```

---

## Fluxo Principal

### 1. Inicialização (`cmd/game/main.go`)

```go
func main() {
    initGLFW()           // Cria contexto de janela
    window := createWindow()
    initOpenGL()         // Inicializa OpenGL
    
    // Carrega modelos 3D
    models = []model.Model{
        model.NewModel("Demogorgon", "assets/dead_by_daylight.glb", ...),
    }
    
    // Compila shaders (programas da GPU)
    shaderProgram = createShaderProgram()
    
    // Game loop
    for !window.ShouldClose() {
        draw()
        window.SwapBuffers()
        glfw.PollEvents()
    }
}
```

### 2. Carregamento de Modelos (`libs/gltfloader/loader.go`)

O `LoadGLB` faz:

1. **Parse do arquivo glTF/GLB** - formato padrão de modelos 3D
2. **Carrega texturas** → `loadTextures`
3. **Processa a árvore de nós** → `processNode` (recursivo)
4. **Converte primitivas para OpenGL** → `loadPrimitive`
   - Cria **VAO/VBO/EBO** (buffers de vértices na GPU)
   - Monta buffer interleaved: `posição(3) + normal(3) + UV(2)`

### 3. Renderização (`draw()`)

```go
func draw() {
    // Matrizes de transformação
    rotMat := matRotateY(float32(angle))      // Rotação
    viewMat := matLookAt(eye, center, up)     // Câmera
    projMat := matPerspective(fov, aspect, near, far)  // Projeção 3D→2D
    
    // Para cada modelo
    for _, entry := range models {
        modelMat := translate * rotate * scale  // Transform do objeto
        
        // Envia para a GPU
        setUniformMat4(shaderProgram, "model", modelMat)
        
        // Desenha
        gl.DrawElements(gl.TRIANGLES, indexCount, ...)
    }
}
```

### 4. Shaders (GLSL)

**Vertex Shader** - roda para cada vértice:
```glsl
gl_Position = projection * view * model * vec4(aPos, 1.0);
vNormal = mat3(transpose(inverse(model))) * aNormal;
```

**Fragment Shader** - roda para cada pixel:
```glsl
// Iluminação difusa simples
float diff = max(dot(norm, light), 0.0);
FragColor = vec4((ambient + diffuse) * color, 1.0);
```

---

## Conceitos para Estudar

### 📐 Matemática 3D (Essencial)

| Conceito | Onde aparece | O que estudar |
|----------|--------------|---------------|
| Vetores | `vecNormalize`, `vecCross` | Operações vetoriais, produto escalar/vetorial |
| Matrizes 4x4 | `mat4`, `matMul` | Multiplicação de matrizes, column-major vs row-major |
| Transformações | `matTranslate`, `matRotateY`, `matScale` | TRS (Translation, Rotation, Scale) |
| Quaternions | `composeTRS` | Rotações sem gimbal lock |
| Projeção | `matPerspective` | Frustum, FOV, near/far planes |
| Câmera | `matLookAt` | View matrix, coordenadas de câmera |

### 🎮 OpenGL

| Conceito | Onde aparece | O que estudar |
|----------|--------------|---------------|
| VAO/VBO/EBO | `loadPrimitive` | Vertex Array Objects, buffers de GPU |
| Shaders GLSL | Vertex/Fragment shaders | Pipeline de renderização |
| Uniforms | `setUniformMat4` | Passagem de dados CPU→GPU |
| Texturas | `uploadImageToGL` | Mapeamento UV, mipmaps |
| Depth Buffer | `gl.Enable(gl.DEPTH_TEST)` | Z-buffer, oclusão |

### 📦 Formato glTF 2.0

| Conceito | Onde aparece |
|----------|--------------|
| Scene Graph | `processNode` |
| Meshes/Primitivas | `loadPrimitive` |
| Materiais PBR | `prim.Material`, `BaseColorFactor` |
| Buffer Views | `loadTextures` |

### 💡 Iluminação

| Conceito | Onde aparece |
|----------|--------------|
| Normais | `generateFlatNormals` |
| Ambient/Diffuse | Fragment shader |
| Flat vs Smooth shading | Cálculo de normais por face |

---

## Estrutura de Dados Chave

```go
// Mesh carregada pronta para OpenGL
type GLTFMesh struct {
    VAO         uint32      // Vertex Array Object
    IndexCount  int32       // Número de índices
    TextureID   uint32      // ID da textura na GPU
    BaseColor   [4]float32  // Cor RGBA do material
    Transform   [16]float32 // Matriz de transformação do nó
}

// Modelo do domínio
type Model struct {
    Scale       Dimensions  // Escala XYZ
    Translation Dimensions  // Posição XYZ
    LoadedModel GLTFModel   // Dados OpenGL
}
```

---

## Recursos de Estudo Recomendados

1. **LearnOpenGL** (https://learnopengl.com) - Tutorial completo de OpenGL moderno
2. **3D Math Primer for Graphics** - Livro de matemática para games
3. **glTF 2.0 Spec** (https://registry.khronos.org/glTF/specs/2.0/glTF-2.0.html)
4. **Essence of Linear Algebra** (3Blue1Brown no YouTube) - Visualização de álgebra linear

---

## 🔦 Task: Implementar Iluminação e Sombreamento

### O que JÁ EXISTE no código

| Componente | Arquivo | Status |
|------------|---------|--------|
| **Normais no buffer** | `loader.go` L300-340 | ✅ `location=1`, stride 32 bytes |
| **Geração de flat normals** | `loader.go` L446-480 | ✅ `generateFlatNormals()` via produto vetorial |
| **Transformação de normais** | `main.go` L118-119 | ✅ `mat3(transpose(inverse(model))) * aNormal` |
| **Luz direcional** | `main.go` L237 | ✅ `lightDir = [-0.3, -0.8, -0.5]` |
| **Ambient + Diffuse** | `main.go` L137-149 | ✅ Básico implementado |

### Shader atual (Fragment)

```glsl
// Ambient fixo
float ambientStrength = 0.2;
vec3 ambient = ambientStrength * vec3(1.0);

// Diffuse (Lambert)
vec3 norm = normalize(vNormal);
vec3 light = normalize(-lightDir);
float diff = max(dot(norm, light), 0.0);
vec3 diffuse = diff * vec3(1.0);

vec3 result = (ambient + diffuse) * color;
```

### O que FALTA implementar

| Feature | Prioridade | Descrição |
|---------|------------|-----------|
| **Specular (Phong/Blinn-Phong)** | Alta | Reflexo brilhante, precisa da posição da câmera |
| **Múltiplas luzes** | Média | Array de luzes, loop no fragment shader |
| **Point Lights** | Média | Atenuação por distância (linear/quadratic) |
| **Spot Lights** | Baixa | Cone de luz com inner/outer cutoff |
| **Shadow Mapping** | Baixa | Depth buffer do ponto de vista da luz |
| **PBR (Physically Based)** | Futuro | Metallic/Roughness do glTF |

---

### 📚 Estudo PRÉ-REQUISITO (entender o código atual)

#### 1. Normais e transformação

**Onde estudar no código:**
- [loader.go](libs/gltfloader/loader.go) linhas 246-260: leitura de normais do glTF
- [loader.go](libs/gltfloader/loader.go) linhas 446-480: `generateFlatNormals()`
- [main.go](cmd/game/main.go) linhas 118-119: transformação no vertex shader

**Conceitos:**
- Por que usar `mat3(transpose(inverse(model)))` para normais?
  - A Normal Matrix corrige distorções causadas por escala não-uniforme
  - Se escalar só em X, a normal não pode ser escalada igual ou fica errada

**Exercício:** Remova o `transpose(inverse(...))` e aplique escala não-uniforme para ver o bug.

#### 2. Iluminação difusa (Lambert)

**Onde estudar no código:**
- [main.go](cmd/game/main.go) linhas 143-147: cálculo do diffuse

**Conceitos:**
```
diffuse = max(dot(N, L), 0.0)
```
- `N` = normal da superfície (normalizada)
- `L` = direção da luz (normalizada)
- `dot(N, L)` = cosseno do ângulo entre eles
- `max(..., 0)` = evita valores negativos (superfície oposta à luz)

#### 3. Direção da luz

**Onde estudar no código:**
- [main.go](cmd/game/main.go) linha 237: `lightDir = [-0.3, -0.8, -0.5]`
- [main.go](cmd/game/main.go) linha 144: `normalize(-lightDir)` inverte para "apontar para a luz"

---

### 📚 Estudo para IMPLEMENTAR (próximos passos)

#### 1. Specular Highlight (Blinn-Phong)

**LearnOpenGL:** https://learnopengl.com/Lighting/Basic-Lighting

**Conceito:**
```glsl
// Blinn-Phong (mais eficiente que Phong puro)
vec3 viewDir = normalize(viewPos - vFragPos);
vec3 halfDir = normalize(lightDir + viewDir);
float spec = pow(max(dot(normal, halfDir), 0.0), shininess);
vec3 specular = specularStrength * spec * lightColor;
```

**O que adicionar no código:**
1. Passar `viewPos` (posição da câmera) como uniform
2. Adicionar `specularStrength` e `shininess` por material
3. Ler `metallicRoughnessTexture` do glTF (opcional)

#### 2. Point Lights com Atenuação

**LearnOpenGL:** https://learnopengl.com/Lighting/Light-casters

**Conceito:**
```glsl
float distance = length(lightPos - fragPos);
float attenuation = 1.0 / (constant + linear * distance + quadratic * distance * distance);
```

**O que adicionar:**
1. Struct `PointLight { vec3 position; vec3 color; float constant, linear, quadratic; }`
2. Array de point lights no shader
3. Loop somando contribuição de cada luz

#### 3. Shadow Mapping

**LearnOpenGL:** https://learnopengl.com/Advanced-Lighting/Shadows/Shadow-Mapping

**Conceito:**
1. Renderizar cena do ponto de vista da luz → depth buffer
2. No fragment shader, comparar profundidade do fragmento com o depth map
3. Se fragmento está "atrás" do depth map → sombra

**Complexidade alta:** requer FBO, nova passada de render, bias para evitar shadow acne.

---

### 🎯 Plano de Implementação Sugerido

```
Fase 1: Specular básico
├── Passar viewPos como uniform
├── Implementar Blinn-Phong no fragment shader
└── Testar com shininess = 32

Fase 2: Estrutura de múltiplas luzes
├── Criar struct Light no GLSL
├── Refatorar para loop de luzes
└── Adicionar 1 point light de teste

Fase 3: Integração com glTF PBR
├── Ler metallicFactor/roughnessFactor do material
├── Mapear roughness → shininess
└── (Opcional) Ler metallic/roughness texture

Fase 4: Shadow Mapping (avançado)
├── Criar FBO para depth map
├── Shadow pass separada
└── Calcular sombras no fragment shader
```
