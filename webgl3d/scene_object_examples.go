package webgl3d

import "github.com/go4orward/gowebgl/common"

func NewSceneObject_3DAxes(wctx *common.WebGLContext, length float32) *SceneObject {
	// This example creates two lines for X (red) and Y (green) axes, with origin at (0,0)
	geometry := NewGeometry() // create an empty geometry
	geometry.SetVertices([][3]float32{{0, 0, 0}, {length, 0, 0}, {0, length, 0}, {0, 0, length}})
	geometry.SetEdges([][]uint32{{0, 1}, {0, 2}, {0, 3}}) // add three edges
	geometry.BuildDataBuffers(true, true, false)          // build data buffers for vertices and edges
	shader := NewShader_3DAxes(wctx)                      // create shader, and set its bindings
	shader.SetThingsToDraw("LINES")                       // let it draw "LINES" later, when Renderer runs
	return NewSceneObject(geometry, nil, shader)          // set up the scene object
}

func NewSceneObject_CylinderWireframe(wctx *common.WebGLContext) *SceneObject {
	// This example creates a hexagon at (0,0), to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works - for faces with more than 3 vertices)
	geometry := NewGeometry_Cylinder(6, 0.5, 1.0, 0, true) // create a cylinder with radius 0.5 and heigt 1.0
	geometry.BuildDataBuffersForWireframe()                // extract wireframe edges from faces
	material := NewMaterial("#888888")                     // create material with yellow color
	shader := NewShader_BasicLight(wctx)                   // create shader, and set its bindings
	shader.SetThingsToDraw("LINES")                        // let it draw "LINES" later, when Renderer runs
	return NewSceneObject(geometry, material, shader)      // set up the scene object
}

func NewSceneObject_SimplePolygon(wctx *common.WebGLContext) *SceneObject {
	geometry := NewGeometry_EmptyExample()
	geometry.BuildDataBuffers(true, false, true)      // build data buffers for vertices and edges
	material := NewMaterial("#888888")                // create material with yellow color
	shader := NewShader_BasicLight(wctx)              // create shader, and set its bindings
	shader.SetThingsToDraw("TRIANGLES")               // let it draw "TRIANGLES" later, when Renderer runs
	return NewSceneObject(geometry, material, shader) // set up the scene object
}
