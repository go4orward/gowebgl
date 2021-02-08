package webgl2d

import "github.com/go4orward/gowebgl/common"

func NewSceneObject_ForAxes(wctx *common.WebGLContext, length float32) *SceneObject {
	// This example creates two lines for X (red) and Y (green) axes, with origin at (0,0)
	geometry := NewGeometry()                                           // create an empty geometry
	geometry.AddVertices([][]float32{{0, 0}, {length, 0}, {0, length}}) // add three vertices
	geometry.AddEdges([][]uint32{{0, 1}, {0, 2}})                       // add two edges
	geometry.BuildDataBuffers(true, true, false)                        // build data buffers for vertices and edges
	geometry.BuildWebGLBuffers(wctx, true, true, false)                 // build WebGL buffers to draw LINES
	shader := NewShader_ForAxes(wctx)                                   // create shader, and set its bindings
	shader.SetBindingToDraw("LINES", geometry.GetDrawBuffer("LINES"), geometry.GetDrawCount("LINES"))
	shader.CheckBindings()
	return NewSceneObject(geometry, nil, shader) // set up the scene object
}

func NewSceneObject_Triangle_Red(wctx *common.WebGLContext) *SceneObject {
	// This example creates a red triangle with radius 0.5 at (0,0)
	geometry := NewGeometry().LoadTriangle(0.5)         // create a triangle with radius 0.5 at (0,0)
	geometry.BuildDataBuffers(true, false, true)        // build data buffers for vertices and faces
	geometry.BuildWebGLBuffers(wctx, true, false, true) // build WebGL buffers to draw TRIANGLES
	shader := NewShader_SimplyRed(wctx)                 // create shader, and set its bindings
	shader.SetBindingForAttribute("pos", "vec2", geometry.GetDrawBuffer("POINTS"), 0, 0)
	shader.SetBindingToDraw("TRIANGLES", geometry.GetDrawBuffer("TRIANGLES"), geometry.GetDrawCount("TRIANGLES"))
	shader.CheckBindings()
	return NewSceneObject(geometry, nil, shader) // set up the scene object
}

func NewSceneObject_Hexagon_Wireframed(wctx *common.WebGLContext) *SceneObject {
	// This example creates a hexagon with given color and radius 0.5 at (0,0), to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works - for faces with more than 3 vertices)
	geometry := NewGeometry().LoadPolygon(6, 0.5, 30)   // create a hexagon with radius 0.5, with 1st vertex at 30 degree from X axis
	geometry.BuildDataBuffersForWireframe()             // extract wireframe edges from faces
	geometry.BuildWebGLBuffers(wctx, true, true, false) // build WebGL buffers to draw LINES
	material := NewMaterial("#ffff00")                  // create material
	shader := NewShader_SingleColor(wctx)               // create shader, and set its bindings
	shader.SetBindingToDraw("LINES", geometry.GetDrawBuffer("LINES"), geometry.GetDrawCount("LINES"))
	shader.CheckBindings()
	return NewSceneObject(geometry, material, shader) // set up the scene object
}
