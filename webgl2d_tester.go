// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl2d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl2d_tester.go'.
	fmt.Println("Hello WebGL!")                       // print in the browser console
	wctx, err := common.NewWebGLContext("wasmcanvas") // canvas_id, interactivity
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	geometry := webgl2d.NewGeometry_Triangle(0.5) // create geometry (a triangle with radius 0.5)
	geometry.BuildDataBuffers(true, false, true)  // build data buffers for vertices and faces
	material := webgl2d.NewMaterial("#bbbbff")    // create material (with light-blue color)
	shader := webgl2d.NewShader_Basic(wctx)       // shader with auto-binded color & PVM matrix
	scnobj := webgl2d.NewSceneObject(geometry, material, shader).Rotate(40)
	scene := webgl2d.NewScene().Add(scnobj)
	camera := webgl2d.NewCamera(wctx.GetWH(), 2.6, 1.0)
	renderer := webgl2d.NewRenderer(wctx) // set up the renderer
	renderer.Clear(camera, "#ffffff")     // prepare to render (clearing to white background)
	renderer.RenderScene(camera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.0)      // render the axes (just for visual reference)
}
