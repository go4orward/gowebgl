// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl3d_example.go'.
	fmt.Println("Hello WebGL!")                       // printed in the browser console
	wctx, err := common.NewWebGLContext("wasmcanvas") // ID of canvas element
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	scene := webgl3d.NewScene("#ffffff")
	if false {
		scene.Add(webgl3d.NewSceneObject_CylinderWireframe(wctx)) // a pre-defined example of SceneObject
		// scene.Add(webgl3d.NewSceneObject_CubeInstances(wctx)) // a pre-defined example of SceneObject
	} else {
		geometry := webgl3d.NewGeometry_CubeWithTexture(1, 1, 1)   // create geometry (a cube of size 1.0)
		geometry.BuildNormalsForFace()                             // calculate normal vectors for each face
		geometry.BuildDataBuffers(true, false, true)               // build data buffers for vertices and faces
		material := common.NewMaterial(wctx, "/assets/gopher.png") // create material (with texture image)
		shader := webgl3d.NewShader_NormalTexture(wctx)            // use the standard NORMAL+TEXTURE shader
		scene.Add(webgl3d.NewSceneObject(geometry, material, nil, nil, shader))
	}
	camera := webgl3d.NewPerspectiveCamera(wctx.GetWH(), 15, 1.0) // FOV default is 15° (in degree)
	camera.SetPose([3]float32{0, 0, 10}, [3]float32{0, 0, 0}, [3]float32{0, 1, 0})
	renderer := webgl3d.NewRenderer(wctx) // set up the renderer
	renderer.Clear(scene)                 // prepare to render (clearing to white background)
	renderer.RenderScene(scene, camera)   // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.0)      // render the axes (just for visual reference)

	if true { // interactive
		fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console
		// add user interactions (with mouse)
		wctx.SetupEventHandlers()
		wctx.RegisterEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
			camera.ShowInfo()
		})
		wctx.RegisterEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
			camera.RotateAroundPoint(10, float32(dxy[0])*0.2, float32(dxy[1])*0.2)
		})
		wctx.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
		})
		wctx.RegisterEventHandlerForWindowResize(func(w int, h int) {
			camera.SetAspectRatio(w, h)
		})
		// add animation
		wctx.SetupAnimationFrame(func(canvas js.Value) {
			renderer.Clear(scene)               // prepare to render (clearing to white background)
			renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(camera, 0.8)    // render the axes (just for visual reference)
			// scene.Get(0).Rotate([3]float32{1, 0, 0}, 1.0)
			scene.Get(0).Rotate([3]float32{0, 1, 0}, 1.0)
			scene.Get(0).Rotate([3]float32{0, 0, 1}, 1.0)
		})
		<-make(chan bool) // wait for events (without exiting)
	}
}
