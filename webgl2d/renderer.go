package webgl2d

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type Renderer struct {
	wctx *common.WebGLContext
	axes *SceneObject
}

func NewRenderer(wctx *common.WebGLContext) *Renderer {
	renderer := Renderer{wctx: wctx, axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *Renderer) Clear(color string) {
	// // Enable the depth test
	// gl.enable(gl.DEPTH_TEST)
	// gl.depthFunc(gl.LEQUAL) // Near things obscure far things
	// // Set the view port
	// // gl.viewport(0, 0, this.viewportWidth, this.viewportHeight);
	// gl.viewport(0, 0, gl.canvas.width, gl.canvas.height)
	// // Clear the color buffer bit
	// gl.clearColor(clear_color[0], clear_color[1], clear_color[2], 1.0)
	// gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()

	rgb, _ := parse_hex_color(color)
	context.Call("clearColor", rgb[0], rgb[1], rgb[2], 1.0) // Set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT)       // Clear the canvas
	// context.Call("enable", constants.DEPTH_TEST)         // Enable the depth test
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *Renderer) RenderAxes(camera *Camera, length float32) {
	if self.axes == nil {
		self.axes = NewSceneObject_ForAxes(self.wctx, length)
	}
	self.RenderSceneObject(self.axes, camera.viewmatrix)
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(camera *Camera, scene *Scene) {
	// Render all the scene objects
	for _, sobj := range scene.Objects {
		modelview := camera.viewmatrix.MultiplyRight(sobj.modelmatrix)
		self.RenderSceneObject(sobj, modelview)
	}
}

// ----------------------------------------------------------------------------
// Rendering SceneObject
// ----------------------------------------------------------------------------

func (self *Renderer) RenderSceneObject(sobj *SceneObject, modelview *geom2d.Matrix3) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	shader := sobj.shader
	if shader == nil {
		shader = sobj.parent_shader
		if shader == nil {
			return errors.New("Failed to RenderSceneObject() : shader not found")
		}
	}
	// 1.
	shader.UseProgram()
	// 2. bind the uniforms of the shader program
	for uname, umap := range shader.uniforms {
		if umap["location"] == nil {
			err := errors.New("Invalid binding : call 'shader.CheckBinding()' before rendering")
			fmt.Println(err.Error())
			return err
		}
		location, dtype := umap["location"].(js.Value), umap["dtype"].(string)
		autobinding, value := umap["autobinding"].(string), umap["value"]
		var err error = nil
		if autobinding != "" {
			err = self.complete_uniform_binding_automatically(location, dtype, autobinding, sobj, modelview)
		} else if value != nil {
			err = self.complete_uniform_binding_with_value(location, dtype, value)
		} else {
			err = errors.New("Invalid binding : uniform '" + uname + "' failed to bind")
		}
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 3. bind the attributes of the shader program
	for aname, amap := range shader.attributes {
		if amap["location"] == nil {
			err := errors.New("Invalid binding : call 'shader.CheckBinding()' before rendering")
			fmt.Println(err.Error())
			return err
		}
		location, dtype := amap["location"].(js.Value), amap["dtype"].(string)
		autobinding, buffer := amap["autobinding"].(string), amap["buffer"]
		var err error = nil
		if autobinding != "" {
			err = self.complete_attribute_binding_automatically(location, dtype, autobinding, sobj)
		} else if buffer != nil {
			stride, offset := amap["stride"].(int), amap["offset"].(int)
			err = self.complete_attribute_binding_with_buffer(location, dtype, buffer, stride, offset)
		} else {
			err = errors.New("Invalid binding : attribute '" + aname + "' failed to bind")
		}
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 4. draw
	for mode, mode_map := range shader.draw_modes {
		// fmt.Println("mode_map : %T %T", mode_map["count"], mode_map["buffer"])
		count := mode_map["count"].(int)
		if count <= 0 {
			fmt.Printf("Invalid BindingsToDraw : %s %v\n", mode, mode_map)
			continue
		}
		switch mode {
		case "POINTS":
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, nil)
			context.Call("drawArrays", constants.POINTS, 0, count)
		case "LINES":
			wgl_buffer := mode_map["buffer"].(js.Value)
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, wgl_buffer)
			context.Call("drawElements", constants.LINES, count, constants.UNSIGNED_INT, 0)
		case "TRIANGLES":
			wgl_buffer := mode_map["buffer"].(js.Value)
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, wgl_buffer)
			context.Call("drawElements", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0)
		default:
		}
	}
	for _, child := range sobj.children {
		modelview := modelview.MultiplyRight(child.modelmatrix)
		self.RenderSceneObject(child, modelview)
	}
	return nil
}

func (self *Renderer) complete_uniform_binding_automatically(location js.Value, dtype string, autobinding string, sobj *SceneObject, modelview *geom2d.Matrix3) error {
	context := self.wctx.GetContext()
	// fmt.Printf("Uniform (%s) : autobinding= '%s'\n", dtype, autobinding)
	switch autobinding {
	case "material.color":
		material := sobj.material
		if material == nil {
			material = sobj.parent_material
			if material == nil {
				return fmt.Errorf("Invalid binding : uniform (%s) - material not found", dtype)
			}
		}
		c := material.color
		switch dtype {
		case "vec3":
			context.Call("uniform3f", location, c[0], c[1], c[2])
			return nil
		case "vec4":
			context.Call("uniform4f", location, c[0], c[1], c[2], c[3])
			return nil
		}
	case "renderer.modelview":
		switch dtype {
		case "mat3":
			// Note that we need Transpose(), since WebGL uses column-major matrix
			e := common.ConvertGoSliceToJsTypedArray(modelview.Transpose().GetElements()) // ModelView matrix, converted to JavaScript 'Float32Array'
			context.Call("uniformMatrix3fv", location, false, e)                          // gl.uniformMatrix3fv(location, transpose, values_array)
			return nil
		}
	}
	return fmt.Errorf("Invalid binding : uniform (%s) failed to bind with '%s'", dtype, autobinding)
}

func (self *Renderer) complete_uniform_binding_with_value(location js.Value, dtype string, value interface{}) error {
	context := self.wctx.GetContext()
	// fmt.Printf("Uniform (%s) : value= %v (%T)\n", dtype, value, value)
	switch dtype {
	case "float":
		context.Call("uniform1f", location, value.(float32))
		return nil
	case "vec2":
		v := value.([]float32)
		context.Call("uniform2f", location, v[0], v[1])
		return nil
	case "vec3":
		v := value.([]float32)
		context.Call("uniform3f", location, v[0], v[1], v[2])
		return nil
	case "vec4":
		v := value.([]float32)
		context.Call("uniform4f", location, v[0], v[1], v[2], v[3])
		return nil
	}
	return fmt.Errorf("Invalid binding : uniform (%s) failed to bind with value %T", dtype, value)
}

func (self *Renderer) complete_attribute_binding_automatically(location js.Value, dtype string, autobinding string, sobj *SceneObject) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	// fmt.Printf("Attribute (%s) : autobinding= '%s'\n", dtype, autobinding)
	switch autobinding {
	case "geometry.coord":
		buffer, stride, offset := sobj.geometry.GetDrawBuffer("POINTS"), 0, 0
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, stride, offset)
		context.Call("enableVertexAttribArray", location)
		return nil
	}
	return fmt.Errorf("Invalid binding : attribute (%s) failed to bind with '%s'", dtype, autobinding)
}

func (self *Renderer) complete_attribute_binding_with_buffer(location js.Value, dtype string, buffer interface{}, stride int, offset int) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	// fmt.Printf("Attribute (%s) : buffer= %v (%T)\n", dtype, value, buffer)
	switch dtype {
	case "vec2":
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer.(js.Value))
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, stride, offset)
		context.Call("enableVertexAttribArray", location)
		return nil
	}
	return fmt.Errorf("Invalid binding : attribute ('%s') failed to bind with buffer %T", dtype, buffer)
}
