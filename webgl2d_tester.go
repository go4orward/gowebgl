// +build js,wasm
package main

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl2d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl2d_tester.go'.
	fmt.Println("Hello WebGL!")                        // print in the browser console
	wctx, err := common.NewWebGLContext("webglcanvas") // canvas_id, interactivity
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
		material := webgl2d.NewMaterial("#bbbbff")          // create material (yellow color)
		shader := webgl2d.NewShader_ModelView(wctx)         // create a shader, and set its bindings
		shader.SetThingsToDraw("TRIANGLES")
		sobj = webgl2d.NewSceneObject(geometry, material, shader)
	}
	scene := webgl2d.NewScene().Add(sobj) // set up the scene, using the SceneObject
	camera := webgl2d.NewCamera()         // set up the camera (by default, centered at (0,0) with zoom 1.0)
	renderer := webgl2d.NewRenderer(wctx) // set up the renderer
	renderer.Clear("#ffffff")             // prepare to render (clearing to white background)
	renderer.RenderScene(camera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 0.8)      // render the axes (just for visual reference)

	if true { // ONLY FOR INTERACTIVE UI
		// add user interactions (with mouse)
		wctx.SetGoCallbackForEventHandling("goProcessEvent") // called from Javascript
		wctx.RegisterEventHandlerForMouseDrag(func(cxy [2]int, sxy [2]int, keystat [4]bool) {
			width := float32(wctx.GetWidth()) * 8
			camera.AddCenter(float32(-cxy[0]+sxy[0])/width, float32(cxy[1]-sxy[1])/width)
		})
		wctx.RegisterEventHandlerForMouseWheel(func(scale float32) {
			scale = float32(math.Sqrt(float64(scale)))
			camera.SetZoom(scale)
		})
		// add animation
		var rotation_angle float32 = 0.0
		wctx.SetGoCallbackForAnimationFrame("goDrawAnimationFrame") // called from Javascript
		wctx.RegisterDrawHandlerForAnimationFrame(func(canvas js.Value) {
			renderer.Clear("#ffffff")           // prepare to render (clearing to white background)
			renderer.RenderScene(camera, scene) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(camera, 0.8)    // render the axes (just for visual reference)
			sobj.Rotate(rotation_angle)
			rotation_angle += 1.0 // animate rotating scene object
		})
		<-make(chan bool) // wait for events from Javascript (without exiting)
	}
}
