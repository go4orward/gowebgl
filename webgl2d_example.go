// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/webgl2d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl2d_example.go'.
	fmt.Println("Hello WebGL!")                        // printed in the browser console
	wctx, err := wcommon.NewWebGLContext("wasmcanvas") // ID of canvas element
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	geometry := webgl2d.NewGeometry_Triangle(0.5)    // create geometry (a triangle with radius 0.5)
	geometry.BuildDataBuffers(true, true, true)      // build data buffers for vertices, edges and faces
	material := wcommon.NewMaterial(wctx, "#bbbbff") // create material (with light-blue color)
	material.SetColorForDrawMode(2, "#0000ff")       // set edge color for the material
	shader := webgl2d.NewShader_MaterialColor(wctx)  // shader with auto-binded color & PVM matrix
	scnobj := webgl2d.NewSceneObject(geometry, material, nil, shader, shader).Rotate(40)
	scene := webgl2d.NewScene("#ffffff").Add(scnobj) // scene holds all the SceneObjects to be rendered
	camera := webgl2d.NewCamera(wctx.GetWH(), 2, 1)  // FOV 2 means range of [-1,+1] in X, ZoomLevel is 1.0
	renderer := webgl2d.NewRenderer(wctx)            // set up the renderer
	renderer.Clear(scene)                            // prepare to render (clearing to white background)
	renderer.RenderScene(scene, camera)              // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.0)                 // render the axes (just for visual reference)
}
