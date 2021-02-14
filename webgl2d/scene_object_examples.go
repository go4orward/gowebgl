package webgl2d

import "github.com/go4orward/gowebgl/common"

func NewSceneObject_2DAxes(wctx *common.WebGLContext, length float32) *SceneObject {
	// This example creates two lines for X (red) and Y (green) axes, with origin at (0,0)
	geometry := NewGeometry()                                            // create an empty geometry
	geometry.SetVertices([][2]float32{{0, 0}, {length, 0}, {0, length}}) // add three vertices
	geometry.SetEdges([][]uint32{{0, 1}, {0, 2}})                        // add two edges
	geometry.BuildDataBuffers(true, true, false)                         // build data buffers for vertices and edges
	shader := NewShader_2DAxes(wctx)                                     // create shader, and set its bindings
	shader.SetThingsToDraw("LINES")                                      // let it draw "LINES" later, when Renderer runs
	return NewSceneObject(geometry, nil, shader)                         // set up the scene object
}

func NewSceneObject_RedTriangle(wctx *common.WebGLContext) *SceneObject {
	// This example creates a red triangle with radius 0.5 at (0,0)
	geometry := NewGeometry_Triangle(0.5)        // create a triangle with radius 0.5 at (0,0)
	geometry.BuildDataBuffers(true, false, true) // build data buffers for vertices and faces
	shader := NewShader_SimplyRed(wctx)          // create shader, and set its bindings
	shader.SetThingsToDraw("TRIANGLES")          // let it draw "TRIANGLES" later, when Renderer runs
	return NewSceneObject(geometry, nil, shader) // set up the scene object
}

func NewSceneObject_HexagonWireframe(wctx *common.WebGLContext) *SceneObject {
	// This example creates a hexagon with given color and radius 0.5 at (0,0), to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works - for faces with more than 3 vertices)
	geometry := NewGeometry_Polygon(6, 0.5, 30)       // create a hexagon with radius 0.5, with 1st vertex at 30 degree from X axis
	geometry.BuildDataBuffersForWireframe()           // extract wireframe edges from faces
	material := NewMaterial("#888888")                // create material
	shader := NewShader_Basic(wctx)                   // create shader, and set its bindings
	shader.SetThingsToDraw("LINES")                   // let it draw "LINES" later, when Renderer runs
	return NewSceneObject(geometry, material, shader) // set up the scene object
}
