package webgl2d

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
)

type Shader struct {
	wctx *common.WebGLContext //

	vshader_code   string   // vertex   shader source code
	fshader_code   string   // fragment shader source code
	vert_shader    js.Value //
	frag_shader    js.Value //
	shader_program js.Value //

	uniforms   map[string]map[string]interface{} // shader uniforms to bind
	attributes map[string]map[string]interface{} // shader attributes to bind
	draw_modes map[string]map[string]interface{} // draw_modes to render ("POINTS", "LINES", "TRIANGLES")
}

func NewShader(wctx *common.WebGLContext, vertex_shader string, fragment_shader string) *Shader {
	shader := Shader{wctx: wctx, vshader_code: vertex_shader, fshader_code: fragment_shader}
	context := shader.wctx.GetContext()
	constants := shader.wctx.GetConstants()
	shader.vert_shader = context.Call("createShader", constants.VERTEX_SHADER) // Create a vertex shader object
	context.Call("shaderSource", shader.vert_shader, shader.vshader_code)      // Attach vertex shader source code
	context.Call("compileShader", shader.vert_shader)                          // Compile the vertex shader
	// if (!gl.getShaderParameter(vertexShader, gl.COMPILE_STATUS)) {
	// 	console.error("GLEX ERROR : AxesRenderer program failure - invalid vertex shader : " + gl.getShaderInfoLog(vertexShader));
	// 	return false;
	// }
	shader.frag_shader = context.Call("createShader", constants.FRAGMENT_SHADER) // Create fragment shader object
	context.Call("shaderSource", shader.frag_shader, shader.fshader_code)        // Attach fragment shader source code
	context.Call("compileShader", shader.frag_shader)                            // Compile the fragmentt shader
	// if (!gl.getShaderParameter(fragmentShader, gl.COMPILE_STATUS)) {
	// 	console.error("GLEX ERROR : AxesRenderer program failure - invalid fragment shader : " + gl.getShaderInfoLog(fragmentShader));
	// 	return false;
	// }
	shader.shader_program = context.Call("createProgram")                   // Create a shader program object to store the combined shader program
	context.Call("attachShader", shader.shader_program, shader.vert_shader) // Attach a vertex shader
	context.Call("attachShader", shader.shader_program, shader.frag_shader) // Attach a fragment shader
	context.Call("linkProgram", shader.shader_program)                      // Link both the programs
	// if (!gl.getProgramParameter(webgl_program, gl.LINK_STATUS)) {
	// 	console.error("GLEX ERROR : AxesRenderer program failure - link failure for shaders : " + gl.getProgramInfoLog(webgl_program));
	// 	return false;
	// }
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	shader.draw_modes = map[string]map[string]interface{}{}
	return &shader
}

// ----------------------------------------------------------------------------
// Bindings for Uniforms (in the shader code)
// ----------------------------------------------------------------------------

func (self *Shader) InitBindingForUniform(name string, dtype string, autobinding string) {
	// Initialize uniform binding with its name, data_type, and auto_binding option.
	// If 'autobinding' is given, then the binding will be attempted automatically.
	//    (examples: "material.color")
	// Otherwise, 'shader.SetBindingForUniform()' has to be called manually.
	self.uniforms[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
}

func (self *Shader) SetBindingForUniform(name string, dtype string, value interface{}) error {
	var err error = nil
	if umap := self.uniforms[name]; umap != nil {
		if umap["dtype"] == dtype {
			umap["value"] = value
		} else {
			err = errors.New(fmt.Sprintf("Setting uniform '%s' failed (invalid type '%s')", name, dtype))
		}
	} else {
		err = errors.New(fmt.Sprintf("Setting uniform '%s' failed (not found)", name))
	}
	return err
}

// ----------------------------------------------------------------------------
// Bindings for Attributes (in the shader code)
// ----------------------------------------------------------------------------

func (self *Shader) InitBindingForAttribute(name string, dtype string, autobinding string) {
	// Initialize attribute binding with its name, data_type, and auto_binding option.
	// If 'autobinding' is given, then the binding will be attempted automatically.
	//    (examples: "geometry.coord")
	// Otherwise, 'shader.SetBindingForUniform()' has to be called manually.
	self.attributes[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
}

func (self *Shader) SetBindingForAttribute(name string, dtype string, buffer interface{}, stride int, offset int) error {
	var err error = nil
	if amap := self.attributes[name]; amap != nil {
		if amap["dtype"] == dtype {
			amap["buffer"] = buffer
			amap["stride"] = stride
			amap["offset"] = offset
		} else {
			err = errors.New(fmt.Sprintf("Shader.SetAttributeBuffer() failed : attribute '%s' invalid type '%s'", name, dtype))
		}
	} else {
		err = errors.New(fmt.Sprintf("Shader.SetUniformValue() failed : attribute '%s' not found", name))
	}
	return err
}

func (self *Shader) SetBindingToDraw(mode string, webgl_buffer js.Value, count int) {
	mode_map := self.draw_modes[mode]
	if mode_map == nil {
		mode_map = map[string]interface{}{}
		self.draw_modes[mode] = mode_map
	}
	mode_map["buffer"] = webgl_buffer
	mode_map["count"] = count
}

func (self *Shader) CheckBindings() {
	// check uniform locations before rendering (since gl.getXXX() is expensive)
	context := self.wctx.GetContext()
	for uname, umap := range self.uniforms {
		location := context.Call("getUniformLocation", self.shader_program, uname)
		umap["location"] = location
		if umap["dtype"] == nil || (umap["autobinding"] == "" && umap["value"] == nil) {
			fmt.Printf("Invalid binding for uniform '%s' : %v \n", uname, umap)
		}
	}
	// check attribute locations
	for aname, amap := range self.attributes {
		location := context.Call("getAttribLocation", self.shader_program, aname)
		amap["location"] = location
		if amap["dtype"] == nil || (amap["autobinding"] == "" && amap["buffer"] == nil) {
			fmt.Printf("Invalid binding for attribute '%s' : %v \n", aname, amap)
		}
	}
}

// ----------------------------------------------------------------------------
// Using Shader Program
// ----------------------------------------------------------------------------

func (self *Shader) UseProgram() {
	context := self.wctx.GetContext()
	context.Call("useProgram", self.shader_program) // Use the combined shader program object
}
