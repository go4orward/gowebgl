package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webglobe"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webglobe_tester.go'
	fmt.Println("Hello WebGL!")                       // printed in the browser console
	wctx, err := common.NewWebGLContext("wasmcanvas") // ID of canvas element
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	scene := webglobe.NewScene(wctx)
	wcamera := webglobe.NewWorldCamera(wctx.GetWH(), 15, 1.0) // FOV default is 15° (in degree)
	wcamera.SetPoseByLonLat(0, 0, 10)                         // longitude 0°, latitude 0°, radius(distance) 10.0

	renderer := webglobe.NewRenderer(wctx) // set up the renderer
	renderer.Clear(wcamera, "#222222")     // prepare to render (clearing to black background)
	renderer.RenderScene(wcamera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(wcamera, 1.2)      // render the axes (just for visual reference)

	if true { // interactive
		fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console
		// add user interactions (with mouse)
		wctx.SetupEventHandlers()
		wctx.RegisterEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
			wcamera.ShowInfo()
		})
		wctx.RegisterEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
			wcamera.RotateAroundGlobe(float32(dxy[0])*0.2, float32(dxy[1])*0.2)
		})
		wctx.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			if keystat[3] { // ZOOM, if SHIFT was pressed
				wcamera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
			}
		})
		wctx.RegisterEventHandlerForWindowResize(func(w int, h int) {
			wcamera.SetAspectRatio(w, h)
		})
		// add animation
		wctx.SetupAnimationFrame(func(canvas js.Value) {
			renderer.Clear(wcamera, "#222222")   // prepare to render (clearing to white background)
			renderer.RenderScene(wcamera, scene) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(wcamera, 1.2)    // render the axes (just for visual reference)
			scene.Globe.Rotate(0.1)
		})
		<-make(chan bool) // wait for events (without exiting)
	}
}
