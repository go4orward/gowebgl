// +build js,wasm
package examples

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/wcommon"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o gowebgl.wasm gowebgl/webgl_example.go'.
	fmt.Println("Hello WebGL!")                        // printed in the browser console
	wctx, err := wcommon.NewWebGLContext("wasmcanvas") // ID of canvas element
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return
	}
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	indices := []uint32{2, 1, 0}
	vertex_shader_code := `
		attribute vec3 xyz;
		void main(void) {
			gl_Position = vec4(xyz, 1.0);
		}`
	fragment_shader_code := `
		void main(void) {
			gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
		}`
	context := wctx.GetContext()
	constants := wctx.GetConstants()

	//// Geometry ////
	// var vertices_array = js.TypedArrayOf(vertices)   // Since js.TypedArrayOf() of Go1.11 is no longer supported,
	// var indices_array = js.TypedArrayOf(indices)     // we have to use js.CopyBytesToJS() instead.
	var vertices_array = wcommon.ConvertGoSliceToJsTypedArray(vertices)
	var indices_array = wcommon.ConvertGoSliceToJsTypedArray(indices)
	vertexBuffer := context.Call("createBuffer", constants.ARRAY_BUFFER)                             // create buffer
	context.Call("bindBuffer", constants.ARRAY_BUFFER, vertexBuffer)                                 // bind the buffer
	context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)        // pass data to buffer
	indexBuffer := context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)                      // create index buffer
	context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, indexBuffer)                          // bind the buffer
	context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW) // pass data to the buffer

	//// Shaders ////
	vertShader := context.Call("createShader", constants.VERTEX_SHADER)   // Create a vertex shader object
	context.Call("shaderSource", vertShader, vertex_shader_code)          // Attach vertex shader source code
	context.Call("compileShader", vertShader)                             // Compile the vertex shader
	fragShader := context.Call("createShader", constants.FRAGMENT_SHADER) // Create fragment shader object
	context.Call("shaderSource", fragShader, fragment_shader_code)        // Attach fragment shader source code
	context.Call("compileShader", fragShader)                             // Compile the fragment shader
	shaderProgram := context.Call("createProgram")                        // Create a shader program to combine the two shaders
	context.Call("attachShader", shaderProgram, vertShader)               // Attach the compiled vertex shader
	context.Call("attachShader", shaderProgram, fragShader)               // Attach the compiled fragment shader
	context.Call("linkProgram", shaderProgram)                            // Make the shader program linked
	context.Call("useProgram", shaderProgram)                             // Let the completed shader program to be used

	//// Attributes ////
	loc := context.Call("getAttribLocation", shaderProgram, "xyz")            // Get the location of attribute 'xyz' in the shader
	context.Call("vertexAttribPointer", loc, 3, constants.FLOAT, false, 0, 0) // Point 'xyz' location to the positions of ARRAY_BUFFER
	context.Call("enableVertexAttribArray", loc)                              // Enable the use of attribute 'xyz' from ARRAY_BUFFER

	//// Draw the scene ////
	context.Call("clearColor", 1.0, 1.0, 1.0, 1.0)    // Set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT) // Clear the canvas
	context.Call("enable", constants.DEPTH_TEST)      // Enable the depth test

	//// Draw the geometry ////
	context.Call("drawElements", constants.TRIANGLES, len(indices), constants.UNSIGNED_SHORT, 0)
}
