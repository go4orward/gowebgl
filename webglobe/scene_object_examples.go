package webglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

func NewSceneObject_Globe(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// This Globe model has texture (without face normals and directional lighting).
	// Since data buffer entries are not repeated for each face, it's more efficient.
	geometry := NewGeometry_Globe(1.0, 36, 18)                 // create geometry (a cube of size 1.0)
	geometry.BuildDataBuffers(true, false, true)               // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.jpg") // create material (yellow color)
	shader := webgl3d.NewShader_BasicTexture(wctx)             // create a shader, and set its bindings
	return webgl3d.NewSceneObject(geometry, material, shader)  // set up the scene object
}

func NewSceneObject_GlobeWithLight(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// This Globe model has texture AND face normals (with directional lighting).
	// Since data buffer entries are repeated, it produces about 6 times bigger data buffer.
	geometry := NewGeometry_Globe(1.0, 36, 18)                 // create geometry (a cube of size 1.0)
	geometry.BuildNormalsPerFace()                             // calculate normal vectors for each face
	geometry.BuildDataBuffers(true, false, true)               // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.jpg") // create material (yellow color)
	shader := webgl3d.NewShader_BasicTextureWithLight(wctx)    // create a shader, and set its bindings
	return webgl3d.NewSceneObject(geometry, material, shader)  // set up the scene object
}
