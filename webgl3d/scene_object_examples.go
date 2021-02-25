package webgl3d

import (
	"math"

	"github.com/go4orward/gowebgl/common"
)

func NewSceneObject_3DAxes(wctx *common.WebGLContext, length float32) *SceneObject {
	// This example creates two lines for X (red) and Y (green) axes, with origin at (0,0)
	geometry := NewGeometry() // create an empty geometry
	geometry.SetVertices([][3]float32{{0, 0, 0}, {length, 0, 0}, {0, length, 0}, {0, 0, length}})
	geometry.SetEdges([][]uint32{{0, 1}, {0, 2}, {0, 3}}) // add three edges
	geometry.BuildDataBuffers(true, true, false)          // build data buffers for vertices and edges
	shader := NewShader_3DAxes(wctx)                      // create shader, and set its bindings
	return NewSceneObject(geometry, nil, shader)          // set up the scene object
}

func NewSceneObject_CylinderWireframe(wctx *common.WebGLContext) *SceneObject {
	// This example creates a cylinder, to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works)
	geometry := NewGeometry_Cylinder(6, 0.5, 1.0, 0, true) // create a cylinder with radius 0.5 and heigt 1.0
	geometry.BuildDataBuffersForWireframe()                // extract wireframe edges from faces
	material := NewMaterial(wctx, "#888888")               // create material with yellow color
	shader := NewShader_NoLight(wctx)                      // create shader, and set its bindings
	return NewSceneObject(geometry, material, shader)      // set up the scene object
}

func NewSceneObject_CubeWithTexture(wctx *common.WebGLContext) *SceneObject {
	geometry := NewGeometry_CubeWithTexture(1.0, 1.0, 1.0)
	geometry.BuildNormalsForFace()
	geometry.BuildDataBuffers(true, false, true)        // build data buffers for vertices and faces
	material := NewMaterial(wctx, "/assets/gopher.png") // create material with a texture image
	shader := NewShader_Basic(wctx)                     // create shader, and set its bindings
	return NewSceneObject(geometry, material, shader)   // set up the scene object
}

func NewSceneObject_CubeInstances(wctx *common.WebGLContext) *SceneObject {
	// This example creates 40,000 instances of a single geometry, each with its own pose (tx, ty)
	geometry := NewGeometry_Cube(0.08, 0.08, 0.08)     // create a cube of size 0.08
	geometry.BuildNormalsForFace()                     // prepare face normal vectors
	geometry.BuildDataBuffers(true, false, true)       //
	material := NewMaterial(wctx, "#888888")           // create material
	shader := NewShader_InstancePoseColor(wctx)        // create shader, and set its bindings
	sobj := NewSceneObject(geometry, material, shader) // set up the scene object
	poses := NewSceneObjectPoses(6, 10*10*10)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				poses.SetPose(i*100+j*10+k, 0, float32(i)/10, float32(j)/10, float32(k)/10) // tx, ty, tz
				ii, jj, kk := math.Abs(float64(i)-5)/5, math.Abs(float64(j)-5)/5, math.Abs(float64(k)-5)/5
				r, g, b := float32(ii), float32(jj), float32(kk)
				poses.SetPose(i*100+j*10+k, 3, r, g, b) // color
			}
		}
	}
	sobj.SetInstancePoses(poses)
	sobj.Translate(-0.5, -0.5, -0.5)
	return sobj
}
