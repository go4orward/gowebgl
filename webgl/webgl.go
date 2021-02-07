// Package webgl implements the simplest library to explain how Go + WASM + WebGL works.
// This code is inspired by https://github.com/bobcob7/wasm-basic-triangle,
// and the issue of js.TypedArrayOf() (of Go1.11) was fixed.

package webgl

import (
	"../common"
)

func GetExampleGeometry(example string) ([]float32, []uint32) {
	var vertices []float32
	var indices []uint32
	switch example {
	case "triangle":
		vertices = []float32{0, 0.7, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
		indices = []uint32{2, 1, 0}
	default: // rectangle
		vertices = []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0, 0.5, 0.5, 0}
		indices = []uint32{2, 1, 0, 2, 0, 3}
	}
	return vertices, indices
}

func GetExampleShaders(example string) (string, string) {
	var vertex_shader_code string
	var fragment_shader_code string
	switch example {
	default: // simple test shaders
		vertex_shader_code = `
		attribute vec3 coordinates;
		void main(void) {
			gl_Position = vec4(coordinates, 1.0);
		}`
		fragment_shader_code = `
		void main(void) {
			gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
		}`
	}
	return vertex_shader_code, fragment_shader_code
}

func Render(wctx *common.WebGLContext, vertices []float32, indices []uint32, vertex_shader_code string, fragment_shader_code string) {
	context := wctx.GetContext()
	constants := wctx.GetConstants()
	//// Geometry ////

	// var vertices_array = js.TypedArrayOf(vertices)   // Since js.TypedArrayOf() of Go1.11 is no longer supported,
	// var indices_array = js.TypedArrayOf(indices)     //
	var vertices_array = common.ConvertGoSliceToJsTypedArray(vertices) // We have to use js.CopyBytesToJS() instead
	var indices_array = common.ConvertGoSliceToJsTypedArray(indices)   //

	vertexBuffer := context.Call("createBuffer", constants.ARRAY_BUFFER)                      // create buffer
	context.Call("bindBuffer", constants.ARRAY_BUFFER, vertexBuffer)                          // bind the buffer
	context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW) // pass data to buffer

	indexBuffer := context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)                      // create index buffer
	context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, indexBuffer)                          // bind the buffer
	context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW) // pass data to the buffer

	//// Shaders ////

	vertShader := context.Call("createShader", constants.VERTEX_SHADER) // Create a vertex shader object
	context.Call("shaderSource", vertShader, vertex_shader_code)        // Attach vertex shader source code
	context.Call("compileShader", vertShader)                           // Compile the vertex shader

	fragShader := context.Call("createShader", constants.FRAGMENT_SHADER) // Create fragment shader object
	context.Call("shaderSource", fragShader, fragment_shader_code)        // Attach fragment shader source code
	context.Call("compileShader", fragShader)                             // Compile the fragmentt shader

	shaderProgram := context.Call("createProgram")          // Create a shader program object to store the combined shader program
	context.Call("attachShader", shaderProgram, vertShader) // Attach a vertex shader
	context.Call("attachShader", shaderProgram, fragShader) // Attach a fragment shader
	context.Call("linkProgram", shaderProgram)              // Link both the programs
	context.Call("useProgram", shaderProgram)               // Use the combined shader program object

	//// Attributes ////

	coord := context.Call("getAttribLocation", shaderProgram, "coordinates")    // Get the attribute location
	context.Call("vertexAttribPointer", coord, 3, constants.FLOAT, false, 0, 0) // Point an attribute to the currently bound VBO
	context.Call("enableVertexAttribArray", coord)                              // Enable the attribute

	//// Draw the scene ////

	context.Call("clearColor", 0.1, 0.1, 0.1, 1.0)                    // Set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT)                 // Clear the canvas
	context.Call("enable", constants.DEPTH_TEST)                      // Enable the depth test
	context.Call("viewport", 0, 0, wctx.GetWidth(), wctx.GetHeight()) // Set the view port

	// Draw the geometry
	context.Call("drawElements", constants.TRIANGLES, len(indices), constants.UNSIGNED_SHORT, 0)

}
