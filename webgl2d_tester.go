// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
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
	scene := webgl2d.NewScene()
	if false {
		scene.Add(webgl2d.NewSceneObject_HexagonWireframe(wctx)) // a pre-defined example of SceneObject
		// scene.Add(webgl2d.NewSceneObject_InstancePoses(wctx)) // a pre-defined example of SceneObject
	} else {
		geometry := webgl2d.NewGeometry_Triangle(0.5) // create geometry (a triangle with radius 0.5)
		geometry.BuildDataBuffers(true, false, true)  // build data buffers for vertices and faces
		material := webgl2d.NewMaterial("#bbbbff")    // create material (with sky-blue color)
		shader := webgl2d.NewShader_Basic(wctx)       // shader with auto-binded color & PVM matrix
		shader.SetThingsToDraw("TRIANGLES")
		scene.Add(webgl2d.NewSceneObject(geometry, material, shader))
	}
	camera := webgl2d.NewCamera(wctx.GetWH(), 2.0, 1.0)
	renderer := webgl2d.NewRenderer(wctx) // set up the renderer
	renderer.Clear(camera, "#ffffff")     // prepare to render (clearing to white background)
	renderer.RenderScene(camera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.0)      // render the axes (just for visual reference)

	if true { // ONLY FOR INTERACTIVE UI
		// add user interactions (with mouse)
		wctx.SetupEventHandlers()
		wctx.RegisterEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
			wxy := camera.UnprojectCanvasToWorld(canvasxy)
			fmt.Printf("canvas (%d %d)  world (%.2f %.2f)\n", canvasxy[0], canvasxy[1], wxy[0], wxy[1])
		})
		wctx.RegisterEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
			camera.ShowInfo()
		})
		wctx.RegisterEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
			wxy := camera.UnprojectCanvasDeltaToWorld(dxy)
			camera.Translate(-wxy[0], -wxy[1])
		})
		wctx.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			oldxy := camera.UnprojectCanvasToWorld(canvasxy)
			camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
			newxy := camera.UnprojectCanvasToWorld(canvasxy)
			delta := geom2d.SubAB(oldxy, newxy)
			camera.Translate(delta[0], delta[1])
		})
		wctx.RegisterEventHandlerForWindowResize(func(w int, h int) {
			camera.SetAspectRatio(w, h)
		})
		// add animation
		wctx.SetupAnimationFrame()
		wctx.RegisterDrawHandlerForAnimationFrame(func(canvas js.Value) {
			renderer.Clear(camera, "#ffffff")   // prepare to render (clearing to white background)
			renderer.RenderScene(camera, scene) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(camera, 0.8)    // render the axes (just for visual reference)
			scene.Get(0).Rotate(1.0)
		})
		<-make(chan bool) // wait for events (without exiting)
	}
}
