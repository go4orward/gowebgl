package wcommon

import (
	"errors"
	"fmt"
	"strings"
	"syscall/js"
)

type Shader struct {
	wctx *WebGLContext //

	vshader_code   string   // vertex   shader source code
	fshader_code   string   // fragment shader source code
	vert_shader    js.Value //
	frag_shader    js.Value //
	shader_program js.Value //
	err            error    //

	uniforms   map[string]map[string]interface{} // shader uniforms to bind
	attributes map[string]map[string]interface{} // shader attributes to bind
}

func NewShader(wctx *WebGLContext, vertex_shader string, fragment_shader string) (*Shader, error) {
	shader := Shader{wctx: wctx, vshader_code: vertex_shader, fshader_code: fragment_shader}
	context := shader.wctx.GetContext()
	constants := shader.wctx.GetConstants()
	shader.vert_shader = context.Call("createShader", constants.VERTEX_SHADER) // Create a vertex shader object
	context.Call("shaderSource", shader.vert_shader, shader.vshader_code)      // Attach vertex shader source code
	context.Call("compileShader", shader.vert_shader)                          // Compile the vertex shader
	if context.Call("getShaderParameter", shader.vert_shader, constants.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(context.Call("getShaderInfoLog", shader.vert_shader).String())
		shader.err = errors.New("VShader failed to compile : " + msg)
		fmt.Println(shader.err.Error())
	}
	shader.frag_shader = context.Call("createShader", constants.FRAGMENT_SHADER) // Create fragment shader object
	context.Call("shaderSource", shader.frag_shader, shader.fshader_code)        // Attach fragment shader source code
	context.Call("compileShader", shader.frag_shader)                            // Compile the fragmentt shader
	if shader.err == nil && context.Call("getShaderParameter", shader.frag_shader, constants.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(context.Call("getShaderInfoLog", shader.frag_shader).String())
		shader.err = errors.New("FShader failed to compile : " + msg)
		fmt.Println(shader.err.Error())
	}
	shader.shader_program = context.Call("createProgram")                   // Create a shader program object to store the combined shader program
	context.Call("attachShader", shader.shader_program, shader.vert_shader) // Attach a vertex shader
	context.Call("attachShader", shader.shader_program, shader.frag_shader) // Attach a fragment shader
	context.Call("linkProgram", shader.shader_program)                      // Link both the programs
	if shader.err == nil && context.Call("getProgramParameter", shader.shader_program, constants.LINK_STATUS).Bool() == false {
		msg := strings.TrimSpace(context.Call("getProgramInfoLog", shader.shader_program).String())
		shader.err = errors.New("ShaderProgram failed to link : " + msg)
		fmt.Println(shader.err.Error())
	}
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	return &shader, shader.err
}

func (self *Shader) GetShaderProgram() js.Value {
	return self.shader_program
}

func (self *Shader) GetUniformBindings() map[string]map[string]interface{} {
	return self.uniforms
}

func (self *Shader) GetAttributeBindings() map[string]map[string]interface{} {
	return self.attributes
}

func (self *Shader) ShowInfo() {
	if self.err == nil {
		fmt.Printf("Shader  OK\n")
	} else {
		fmt.Printf("Shader  with Error - %s\n", self.err.Error())
	}
	for uname, umap := range self.uniforms {
		fmt.Printf("    Uniform   %-10s: %v\n", uname, umap)
	}
	for aname, amap := range self.attributes {
		fmt.Printf("    Attribute %-10s: %v\n", aname, amap)
	}
}

// ----------------------------------------------------------------------------
// Bindings
// ----------------------------------------------------------------------------

func (self *Shader) SetBindingForUniform(name string, dtype string, option interface{}) {
	// Set uniform binding with its name, data_type, and AUTO/MANUAL option.
	switch option.(type) {
	case string: // let Renderer bind the uniform variable automatically
		autobinding := option.(string)
		autobinding_split := strings.Split(option.(string), ":")
		autobinding0 := autobinding_split[0] // "material.texture:0" (with texture UNIT value)
		switch autobinding0 {
		case "lighting.dlight": // [mat3](3D) directional light information with (direction[3], color[3], ambient[3])
		case "material.color": //  [vec3] uniform color taken from Material
		case "material.texture": // [sampler2D] texture sampler(unit), like "material.texture:0"
		case "renderer.aspect": // AspectRatio of camera, Width : Height
		case "renderer.pvm": //  [mat3](2D) or [mat4](3D) (Proj * View * Model) matrix
		case "renderer.proj": // [mat3](2D) or [mat4](3D) (Projection) matrix
		case "renderer.vwmd": // [mat3](2D) or [mat4](3D) (View * Model) matrix
		default:
			fmt.Printf("Failed to SetBindingForUniform('%s') : unknown autobinding '%s'\n", name, autobinding)
			return
		}
		self.uniforms[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
	case []float32: // let Renderer set the uniform manually, with the given values
		self.uniforms[name] = map[string]interface{}{"dtype": dtype, "value": option.([]float32)}
	default:
		fmt.Printf("Failed to SetBindingForUniform('%s') : invalid option %v\n", name, option)
		return
	}
}

func (self *Shader) SetBindingForAttribute(name string, dtype string, autobinding string) {
	// Set attribute binding with its name, data_type, with AUTO_BINDING option.
	autobinding_split := strings.Split(autobinding, ":")
	autobinding0 := autobinding_split[0]
	switch autobinding0 {
	case "geometry.coords": // point coordinates
	case "geometry.textuv": // texture UV coordinates
	case "geometry.normal": // (3D only) normal vector
	case "instance.pose": // instance pose, like "instance.pose:<stride>:<offset>"
		if len(autobinding_split) != 3 {
			fmt.Printf("Failed to SetBindingForAttribute('%s') : try 'instance.pose:<stride>:<offset>'\n", name)
			return
		}
	default:
		fmt.Printf("Failed to SetBindingForAttribute('%s') : invalid autobinding '%s'\n", name, autobinding)
		return
	}
	self.attributes[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
}

func (self *Shader) CheckBindings() {
	// check uniform locations before rendering (since gl.getXXX() is expensive)
	context := self.wctx.GetContext()
	if self.err != nil {
		fmt.Printf("ShaderProgram is not ready for CheckBindings()\n")
		return
	}
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
