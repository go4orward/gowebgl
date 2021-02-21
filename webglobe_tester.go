package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
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
	fmt.Println("Please wait while world image is loaded.") // printed in the browser console
	scnobj := webglobe.NewSceneObject_Globe(wctx)
	scene := webgl3d.NewScene().Add(scnobj)
	camera := webgl3d.NewPerspectiveCamera(wctx.GetWH(), 15, 1.0) // FOV default is 15 in degree
	camera.SetPose([3]float32{0, 0, 10}, [3]float32{0, 0, 0}, [3]float32{0, 1, 0})
	renderer := webgl3d.NewRenderer(wctx) // set up the renderer
	renderer.Clear(camera, "#000000")     // prepare to render (clearing to black background)
	renderer.RenderScene(camera, scene)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.2)      // render the axes (just for visual reference)
	// scene.Get(0).ShowInfo()

	if true { // interactive
		// add user interactions (with mouse)
		wctx.SetupEventHandlers()
		wctx.RegisterEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
			camera.ShowInfo()
		})
		wctx.RegisterEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
			camera.RotateAroundPoint(10, float32(dxy[0])*0.2, float32(dxy[1])*0.2)
			if !keystat[3] {
				camera.RotateByRoll(0)
			}
		})
		wctx.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			if keystat[3] { // ZOOM, if SHIFT was pressed
				camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
			}
		})
		wctx.RegisterEventHandlerForWindowResize(func(w int, h int) {
			camera.SetAspectRatio(w, h)
		})
		// add animation
		wctx.SetupAnimationFrame(func(canvas js.Value) {
			renderer.Clear(camera, "#000000")   // prepare to render (clearing to white background)
			renderer.RenderScene(camera, scene) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(camera, 1.2)    // render the axes (just for visual reference)
			// scene.Get(0).Rotate([3]float32{0, 0, 1}, 1.0)
		})
		<-make(chan bool) // wait for events (without exiting)
	}
}
