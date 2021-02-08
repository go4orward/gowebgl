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
	fmt.Println("Hello WebGL!") // print in the browser console
	wctx, err := common.NewWebGLContext("webglcanvas")
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	var sobj *webgl2d.SceneObject
	if false {
		sobj = webgl2d.NewSceneObject_Hexagon_Wireframed(wctx) // One of pre-defined examples of SceneObject
	} else {
		geometry := webgl2d.NewGeometry().LoadTriangle(0.5) // create geometry (a triangle with radius 0.5)
		geometry.BuildDataBuffers(true, false, true)        // build data buffers for vertices and faces
		geometry.BuildWebGLBuffers(wctx, true, false, true) // build WebGL buffers to draw POINTS, LINES, TRIANGLES
		material := webgl2d.NewMaterial("#cccc00")          // create material (yellow color)
		shader := webgl2d.NewShader_ModelView(wctx)         // create a shader, and set its bindings
		shader.SetBindingToDraw("TRIANGLES", geometry.GetDrawBuffer("TRIANGLES"), geometry.GetDrawCount("TRIANGLES"))
		shader.CheckBindings()
		sobj = webgl2d.NewSceneObject(geometry, material, shader)
	}
	scene := webgl2d.NewScene().Add(sobj) // set up the scene, using the SceneObject
	camera := webgl2d.NewCamera(wctx)     // set up the camera (by default, looking at (0,0) with size 2.0)
	renderer := webgl2d.NewRenderer(wctx) // set up the renderer
	renderer.Clear("#ffffff")             // prepare to render (clearing white background)
	renderer.RenderScene(camera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 0.8)      // render the axes (just for visual reference)
	camera.ShowInfo()
}
