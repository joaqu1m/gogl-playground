package gltfloader

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/libs/logger"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// GLTFMesh contém os dados OpenGL prontos para renderizar.
type GLTFMesh struct {
	Name        string
	VAO         uint32
	VertexCount int32
	IndexCount  int32
	HasIndices  bool
	TextureID   uint32
	HasTexture  bool
	BaseColor   [4]float32
	Transform   [16]float32 // Node world transform, column-major
}

// GLTFModel agrupa todas as meshes carregadas de um arquivo glTF/GLB.
type GLTFModel struct {
	Meshes []*GLTFMesh
}

// LoadGLB carrega um arquivo .glb/.gltf e cria os recursos OpenGL.
// As posições são carregadas cruas, sem normalização. Os transforms dos nós
// da scene graph são armazenados em cada GLTFMesh.Transform.
func LoadGLB(filepath string) (*GLTFModel, error) {
	doc, err := gltf.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("gltfloader: falha ao abrir %q: %w", filepath, err)
	}

	textures, err := loadTextures(doc)
	if err != nil {
		return nil, fmt.Errorf("gltfloader: falha ao carregar texturas: %w", err)
	}

	model := &GLTFModel{}

	if len(doc.Scenes) > 0 {
		// Percorre a scene graph a partir da cena ativa
		sceneIdx := 0
		if doc.Scene != nil {
			sceneIdx = *doc.Scene
		}
		scene := doc.Scenes[sceneIdx]
		for _, nodeIdx := range scene.Nodes {
			if err := processNode(doc, nodeIdx, mat4fIdentity(), model, textures); err != nil {
				return nil, err
			}
		}
	} else {
		// Fallback: sem cenas definidas, carrega todas as meshes com transform identidade
		for _, mesh := range doc.Meshes {
			for _, prim := range mesh.Primitives {
				glMesh, err := loadPrimitive(doc, prim, textures)
				if err != nil {
					return nil, fmt.Errorf("gltfloader: falha ao carregar primitiva de %q: %w", mesh.Name, err)
				}
				glMesh.Transform = mat4fIdentity()
				model.Meshes = append(model.Meshes, glMesh)
			}
		}
	}

	if len(model.Meshes) == 0 {
		return nil, fmt.Errorf("gltfloader: nenhuma mesh encontrada em %q", filepath)
	}

	return model, nil
}

// processNode percorre recursivamente a árvore de nós, acumulando transforms.
func processNode(doc *gltf.Document, nodeIdx int, parentTransform [16]float32, model *GLTFModel, textures map[int]uint32) error {
	if nodeIdx < 0 || nodeIdx >= len(doc.Nodes) {
		return fmt.Errorf("gltfloader: node index %d fora do range", nodeIdx)
	}
	node := doc.Nodes[nodeIdx]

	localTransform := nodeLocalTransform(node)
	worldTransform := mat4fMul(parentTransform, localTransform)

	if node.Mesh != nil {
		meshIdx := *node.Mesh
		if meshIdx < 0 || meshIdx >= len(doc.Meshes) {
			return fmt.Errorf("gltfloader: mesh index %d fora do range", meshIdx)
		}
		mesh := doc.Meshes[meshIdx]
		for _, prim := range mesh.Primitives {
			glMesh, err := loadPrimitive(doc, prim, textures)
			if err != nil {
				return fmt.Errorf("gltfloader: falha ao carregar primitiva de %q: %w", mesh.Name, err)
			}
			glMesh.Name = mesh.Name
			glMesh.Transform = worldTransform
			model.Meshes = append(model.Meshes, glMesh)

			logger.Infof("mesh %q: node=%q transform=[%.3f, %.3f, %.3f, %.3f | %.3f, %.3f, %.3f, %.3f | %.3f, %.3f, %.3f, %.3f | %.3f, %.3f, %.3f, %.3f]",
				mesh.Name, node.Name,
				worldTransform[0], worldTransform[1], worldTransform[2], worldTransform[3],
				worldTransform[4], worldTransform[5], worldTransform[6], worldTransform[7],
				worldTransform[8], worldTransform[9], worldTransform[10], worldTransform[11],
				worldTransform[12], worldTransform[13], worldTransform[14], worldTransform[15],
			)
		}
	}

	for _, childIdx := range node.Children {
		if err := processNode(doc, childIdx, worldTransform, model, textures); err != nil {
			return err
		}
	}

	return nil
}

// nodeLocalTransform calcula a matriz de transformação local de um nó.
func nodeLocalTransform(node *gltf.Node) [16]float32 {
	// A biblioteca qmuntal/gltf inicializa Matrix com DefaultMatrix (identidade)
	// no UnmarshalJSON, mesmo que o JSON não contenha "matrix". Então comparamos
	// contra DefaultMatrix: se for diferente, o nó definiu uma matriz explícita.
	if node.Matrix != gltf.DefaultMatrix {
		var m [16]float32
		for i, v := range node.Matrix {
			m[i] = float32(v)
		}
		return m
	}

	// Caso contrário, compor T * R * S
	t := node.TranslationOrDefault()
	r := node.RotationOrDefault()
	s := node.ScaleOrDefault()
	return composeTRS(
		[3]float32{float32(t[0]), float32(t[1]), float32(t[2])},
		[4]float32{float32(r[0]), float32(r[1]), float32(r[2]), float32(r[3])},
		[3]float32{float32(s[0]), float32(s[1]), float32(s[2])},
	)
}

// composeTRS constrói uma matriz 4x4 column-major a partir de translation,
// rotation (quaternion xyzw) e scale. Resultado = T * R * S.
func composeTRS(t [3]float32, q [4]float32, s [3]float32) [16]float32 {
	x, y, z, w := q[0], q[1], q[2], q[3]
	xx := x * x
	yy := y * y
	zz := z * z
	xy := x * y
	xz := x * z
	yz := y * z
	wx := w * x
	wy := w * y
	wz := w * z

	return [16]float32{
		// Column 0
		(1 - 2*(yy+zz)) * s[0],
		(2 * (xy + wz)) * s[0],
		(2 * (xz - wy)) * s[0],
		0,
		// Column 1
		(2 * (xy - wz)) * s[1],
		(1 - 2*(xx+zz)) * s[1],
		(2 * (yz + wx)) * s[1],
		0,
		// Column 2
		(2 * (xz + wy)) * s[2],
		(2 * (yz - wx)) * s[2],
		(1 - 2*(xx+yy)) * s[2],
		0,
		// Column 3
		t[0], t[1], t[2], 1,
	}
}

// ---- Funções auxiliares de matriz [16]float32 (column-major) ----

func mat4fIdentity() [16]float32 {
	return [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func mat4fMul(a, b [16]float32) [16]float32 {
	var r [16]float32
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

// loadPrimitive converte uma primitiva glTF em VAO/VBO/EBO do OpenGL.
// As posições são carregadas cruas, sem normalização.
func loadPrimitive(doc *gltf.Document, prim *gltf.Primitive, textures map[int]uint32) (*GLTFMesh, error) {
	// ---- Lê posições (obrigatório) ----
	posAccessorIdx, ok := prim.Attributes[gltf.POSITION]
	if !ok {
		return nil, fmt.Errorf("primitiva sem POSITION")
	}
	posData, err := modeler.ReadPosition(doc, doc.Accessors[posAccessorIdx], nil)
	if err != nil {
		return nil, fmt.Errorf("erro lendo posições: %w", err)
	}

	// Log vertex bounds para debug
	if len(posData) > 0 {
		minP := posData[0]
		maxP := posData[0]
		for _, p := range posData[1:] {
			for axis := 0; axis < 3; axis++ {
				if p[axis] < minP[axis] {
					minP[axis] = p[axis]
				}
				if p[axis] > maxP[axis] {
					maxP[axis] = p[axis]
				}
			}
		}
		logger.Infof("  primitive %d verts: bounds min=(%.4f, %.4f, %.4f) max=(%.4f, %.4f, %.4f)",
			len(posData),
			minP[0], minP[1], minP[2],
			maxP[0], maxP[1], maxP[2],
		)
	}

	// ---- Lê normais (opcional) ----
	var normalData [][3]float32
	if normIdx, ok := prim.Attributes[gltf.NORMAL]; ok {
		normalData, err = modeler.ReadNormal(doc, doc.Accessors[normIdx], nil)
		if err != nil {
			normalData = nil // fallback: calcula depois
		}
	}

	// ---- Lê UVs (opcional) ----
	var uvData [][2]float32
	if uvIdx, ok := prim.Attributes[gltf.TEXCOORD_0]; ok {
		uvData, _ = modeler.ReadTextureCoord(doc, doc.Accessors[uvIdx], nil)
	}

	// ---- Lê índices (opcional) ----
	var indices []uint32
	if prim.Indices != nil {
		indData, err := modeler.ReadIndices(doc, doc.Accessors[*prim.Indices], nil)
		if err != nil {
			return nil, fmt.Errorf("erro lendo índices: %w", err)
		}
		indices = indData
	}

	// ---- Material ----
	baseColor := [4]float32{0.8, 0.8, 0.8, 1.0}
	var texID uint32
	hasTexture := false

	if prim.Material != nil {
		mat := doc.Materials[*prim.Material]
		if mat.PBRMetallicRoughness != nil {
			pbr := mat.PBRMetallicRoughness
			bc := pbr.BaseColorFactorOrDefault()
			baseColor = [4]float32{float32(bc[0]), float32(bc[1]), float32(bc[2]), float32(bc[3])}

			if pbr.BaseColorTexture != nil {
				texIdx := pbr.BaseColorTexture.Index
				if id, ok := textures[texIdx]; ok {
					texID = id
					hasTexture = true
				}
			}
		}
	}

	// ---- Gera normais se não existirem ----
	if normalData == nil {
		normalData = generateFlatNormals(posData, indices)
	}

	// ---- Monta buffer interleaved: pos(3) + normal(3) + uv(2) = 8 floats ----
	const floatsPerVert = 8
	vertCount := len(posData)
	buf := make([]float32, 0, vertCount*floatsPerVert)

	for i := 0; i < vertCount; i++ {
		// pos
		buf = append(buf, posData[i][0], posData[i][1], posData[i][2])
		// normal
		if i < len(normalData) {
			buf = append(buf, normalData[i][0], normalData[i][1], normalData[i][2])
		} else {
			buf = append(buf, 0, 1, 0)
		}
		// uv
		if uvData != nil && i < len(uvData) {
			buf = append(buf, uvData[i][0], uvData[i][1])
		} else {
			buf = append(buf, 0, 0)
		}
	}

	// ---- Cria VAO/VBO/EBO ----
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(buf)*4, gl.Ptr(buf), gl.STATIC_DRAW)

	const stride = floatsPerVert * 4 // 32 bytes

	// location 0: posição
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// location 1: normal
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// location 2: texcoord
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	glMesh := &GLTFMesh{
		VAO:        vao,
		HasTexture: hasTexture,
		TextureID:  texID,
		BaseColor:  baseColor,
	}

	if len(indices) > 0 {
		var ebo uint32
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
		glMesh.HasIndices = true
		glMesh.IndexCount = int32(len(indices))
	} else {
		glMesh.VertexCount = int32(vertCount)
	}

	gl.BindVertexArray(0)

	return glMesh, nil
}

// loadTextures carrega todas as texturas do documento glTF e retorna um mapa texIdx -> OpenGL texture ID.
func loadTextures(doc *gltf.Document) (map[int]uint32, error) {
	result := make(map[int]uint32)

	for i, tex := range doc.Textures {
		if tex.Source == nil {
			continue
		}
		imgIdx := *tex.Source
		if imgIdx >= len(doc.Images) {
			continue
		}
		img := doc.Images[imgIdx]

		var imgBytes []byte

		if img.BufferView != nil {
			// Imagem embedded no buffer
			bv := doc.BufferViews[*img.BufferView]
			buf := doc.Buffers[bv.Buffer]
			imgBytes = buf.Data[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]
		} else if img.IsEmbeddedResource() {
			// Imagem embedded como data URI
			data, err := img.MarshalData()
			if err != nil {
				continue
			}
			imgBytes = data
		} else {
			// Imagem externa (URI) não suportada neste loader
			continue
		}

		texID, err := uploadImageToGL(imgBytes)
		if err != nil {
			continue
		}
		result[i] = texID
	}

	return result, nil
}

// uploadImageToGL decodifica bytes de imagem e sobe como textura OpenGL.
func uploadImageToGL(data []byte) (uint32, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	w := int32(rgba.Bounds().Dx())
	h := int32(rgba.Bounds().Dy())

	// Não fazemos flip vertical: glTexImage2D mapeia o primeiro pixel para texcoord (0,0),
	// e glTF UV (0,0) é o topo-esquerda da imagem, que coincide com o primeiro pixel
	// decodificado de PNG/JPEG. As convenções se cancelam.

	var texID uint32
	gl.GenTextures(1, &texID)
	gl.BindTexture(gl.TEXTURE_2D, texID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texID, nil
}

// generateFlatNormals calcula normais por face (flat shading).
func generateFlatNormals(positions [][3]float32, indices []uint32) [][3]float32 {
	normals := make([][3]float32, len(positions))

	processTriangle := func(i0, i1, i2 int) {
		p0, p1, p2 := positions[i0], positions[i1], positions[i2]
		e1 := [3]float32{p1[0] - p0[0], p1[1] - p0[1], p1[2] - p0[2]}
		e2 := [3]float32{p2[0] - p0[0], p2[1] - p0[1], p2[2] - p0[2]}
		n := [3]float32{
			e1[1]*e2[2] - e1[2]*e2[1],
			e1[2]*e2[0] - e1[0]*e2[2],
			e1[0]*e2[1] - e1[1]*e2[0],
		}
		l := float32(math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])))
		if l > 0 {
			n[0] /= l
			n[1] /= l
			n[2] /= l
		}
		normals[i0] = n
		normals[i1] = n
		normals[i2] = n
	}

	if len(indices) > 0 {
		for i := 0; i+2 < len(indices); i += 3 {
			processTriangle(int(indices[i]), int(indices[i+1]), int(indices[i+2]))
		}
	} else {
		for i := 0; i+2 < len(positions); i += 3 {
			processTriangle(i, i+1, i+2)
		}
	}

	return normals
}
