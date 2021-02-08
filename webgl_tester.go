// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"

	"./common"
	"./webgl"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl_tester.go'
	fmt.Println("Hello WebGL!") // print in the browser console
	wctx, err := common.NewWebGLContext("webglcanvas")
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
	} else {
		vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
		indices := []uint32{2, 1, 0}
		v_shader, f_shader := webgl.GetExampleShaders("default")
		webgl.Render(wctx, vertices, indices, v_shader, f_shader)
	}
}
