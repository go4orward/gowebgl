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
	scene.Add(webgl2d.NewSceneObject_RectInstances(wctx)) // a pre-defined example of SceneObject
	camera := webgl2d.NewCamera(wctx.GetWH(), 2.6, 1.0)
	bbox, size, center := scene.GetBBoxSizeCenter(true)
	camera.SetFov(size[0]*1.1).SetPose(center[0], center[1], 0.0).SetTranslationBoundingBox(bbox)
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
			wdxy := camera.UnprojectCanvasDeltaToWorld(dxy)
			camera.Translate(-wdxy[0], -wdxy[1])
		})
		wctx.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			if keystat[1] { // ZOOM
				oldxy := camera.UnprojectCanvasToWorld(canvasxy)
				camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
				newxy := camera.UnprojectCanvasToWorld(canvasxy)
				delta := geom2d.SubAB(newxy, oldxy)
				camera.Translate(-delta[0], -delta[1])
			} else { // SCROLL
				wdxy := camera.UnprojectCanvasDeltaToWorld([2]int{0, int(scale)})
				camera.Translate(0.0, wdxy[1])
			}
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
		})
		<-make(chan bool) // wait for events (without exiting)
	}
}
