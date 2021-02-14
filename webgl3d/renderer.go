package webgl3d

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
	"github.com/go4orward/gowebgl/common/geom3d"
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

func (self *Renderer) Clear(camera *Camera, color string) {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()

	// context.Call("viewport", 0, 0, camera.wh[0], camera.wh[1]) // (LowerLeft.x, LowerLeft.y, width, height)
	// (if 'viewport' is not updated, rendering may blur after window.resize)

	rgb, _ := common.ParseHexColor(color)
	context.Call("clearColor", rgb[0], rgb[1], rgb[2], 1.0) // set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT)       // clear the canvas

	context.Call("enable", constants.DEPTH_TEST)      // Enable depth test
	context.Call("depthFunc", constants.LEQUAL)       // Near things obscure far things
	context.Call("clearColor", 0, 0, 0, 1.0)          // set clearing color to all 0
	context.Call("clear", constants.DEPTH_BUFFER_BIT) // clear the depth_buffer
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *Renderer) RenderAxes(camera *Camera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	if self.axes == nil {
		self.axes = NewSceneObject_3DAxes(self.wctx, length)
	}
	self.RenderSceneObject(self.axes, camera.projection.GetMatrix(), &camera.viewmatrix)
	// camera.TestDataBuffer(self.axes.geometry.data_buffer_vpoints, self.axes.geometry.vpoint_info[0])
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(camera *Camera, scene *Scene) {
	// Render all the SceneObjects in the Scene
	for _, sobj := range scene.objects {
		viewmodel := camera.viewmatrix.MultiplyRight(&sobj.modelmatrix)
		self.RenderSceneObject(sobj, camera.projection.GetMatrix(), viewmodel)
	}
}

// ----------------------------------------------------------------------------
// Rendering SceneObject
// ----------------------------------------------------------------------------

func (self *Renderer) RenderSceneObject(sobj *SceneObject, proj *geom3d.Matrix4, viewm *geom3d.Matrix4) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	// 1. If necessary, then WebGLBuffers for the SceneObject's Geometry
	if sobj.geometry.IsDataBufferReady() == false {
		return errors.New("Failed to RenderSceneObject() : empty geometry data buffer")
	}
	if sobj.geometry.IsWebGLBufferReady() == false {
		sobj.geometry.build_webgl_buffers(self.wctx, true, true, true)
	}
	// 2. Decide which Shader to use
	shader := sobj.shader
	if shader == nil {
		shader = sobj.parent_shader
		if shader == nil {
			return errors.New("Failed to RenderSceneObject() : shader not found")
		}
	}
	context.Call("useProgram", shader.GetShaderProgram())
	// 3. bind the uniforms of the shader program
	for uname, umap := range shader.GetUniformBindings() {
		if umap["location"] == nil {
			err := errors.New("Invalid binding : call 'shader.CheckBinding()' before rendering")
			fmt.Println(err.Error())
			return err
		}
		location, dtype := umap["location"].(js.Value), umap["dtype"].(string)
		autobinding, value := umap["autobinding"].(string), umap["value"]
		var err error = nil
		if autobinding != "" {
			err = self.complete_uniform_binding_automatically(location, dtype, autobinding, sobj, proj, viewm)
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
	// 4. bind the attributes of the shader program
	for aname, amap := range shader.GetAttributeBindings() {
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
	// 5. draw
	for _, draw_mode := range shader.GetThingsToDraw() {
		// Note that ARRAY_BUFFER was binded already in the previous step (during attribute binding)
		switch draw_mode {
		case "POINTS":
			_, count, _ := sobj.geometry.GetWebGLBuffer("POINTS")
			context.Call("drawArrays", constants.POINTS, 0, count) // (mode, first, count)
			// fmt.Printf("Renderer: drawArrays POINTS %d\n", count)
		case "LINES":
			buffer, count, _ := sobj.geometry.GetWebGLBuffer("LINES")
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
			context.Call("drawElements", constants.LINES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
			// fmt.Printf("Renderer: drawElements LINES %d (%v)\n", count, pinfo)
		case "TRIANGLES":
			buffer, count, _ := sobj.geometry.GetWebGLBuffer("TRIANGLES")
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
			context.Call("drawElements", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
			// fmt.Printf("Renderer: drawElements TRIANGLES %d (%v)\n", count, pinfo)
		default:
			err := fmt.Errorf("Unknown mode to draw : %s\n", draw_mode)
			fmt.Printf(err.Error())
			return err
		}
	}
	// 6. render all the children
	for _, child := range sobj.children {
		viewm := viewm.MultiplyRight(&child.modelmatrix)
		self.RenderSceneObject(child, proj, viewm)
	}
	return nil
}

func (self *Renderer) complete_uniform_binding_automatically(location js.Value, dtype string, autobinding string,
	sobj *SceneObject, proj *geom3d.Matrix4, viewm *geom3d.Matrix4) error {
	context := self.wctx.GetContext()
	// fmt.Printf("Uniform (%s) : autobinding= '%s'\n", dtype, autobinding)
	switch autobinding {
	case "renderer.proj": // mat4
		e := (*proj.GetElements())[:]
		m := common.ConvertGoSliceToJsTypedArray(e)          // Projection matrix, converted to JavaScript 'Float32Array'
		context.Call("uniformMatrix4fv", location, false, m) // gl.uniformMatrix4fv(location, transpose, values_array)
		return nil
	case "renderer.vmod": // mat4
		e := (*viewm.GetElements())[:]
		m := common.ConvertGoSliceToJsTypedArray(e)          // View * Models matrix, converted to JavaScript 'Float32Array'
		context.Call("uniformMatrix4fv", location, false, m) // gl.uniformMatrix4fv(location, transpose, values_array)
		return nil
	case "renderer.pvm": // mat4
		pvm := proj.MultiplyRight(viewm) // (Proj * View * Models) matrix
		e := (*pvm.GetElements())[:]
		m := common.ConvertGoSliceToJsTypedArray(e)          // P*V*M matrix, converted to JavaScript 'Float32Array'
		context.Call("uniformMatrix4fv", location, false, m) // gl.uniformMatrix4fv(location, transpose, values_array)
		return nil
	case "material.color":
		c := [4]float32{1, 1, 1, 1}
		if sobj.material != nil {
			c = sobj.material.color
		} else if sobj.parent_material != nil {
			c = sobj.parent_material.color
		}
		switch dtype {
		case "vec3":
			context.Call("uniform3f", location, c[0], c[1], c[2])
			return nil
		case "vec4":
			context.Call("uniform4f", location, c[0], c[1], c[2], c[3])
			return nil
		}
	case "lighting.dlight": // mat3
		dlight := geom2d.NewMatrix3().Set(0, 1, 0, 0, 1, 0, 1, 1, 0) // directional light (in camera space)
		e := (*dlight.GetElements())[:]                              // (direction[3] + intensity[3] + ambient[3])
		m := common.ConvertGoSliceToJsTypedArray(e)                  // converted to JavaScript 'Float32Array'
		context.Call("uniformMatrix3fv", location, false, m)         // gl.uniformMatrix4fv(location, transpose, values_array)
		return nil
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
	buffer, _, pinfo := sobj.geometry.GetWebGLBuffer("POINTS")
	switch autobinding {
	case "geometry.coords":
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 3, constants.FLOAT, false, pinfo[0]*4, pinfo[1]*4)
		context.Call("enableVertexAttribArray", location)
		return nil
	case "geometry.textuv":
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, pinfo[0]*4, pinfo[2]*4)
		context.Call("enableVertexAttribArray", location)
		if pinfo[1] == pinfo[2] {
			fmt.Printf("Renderer Warning : Texture UV coordinates not found (pinfo=%v)\n", pinfo)
		}
		return nil
	case "geometry.normal":
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 3, constants.FLOAT, false, pinfo[0]*4, pinfo[3]*4)
		context.Call("enableVertexAttribArray", location)
		if pinfo[1] == pinfo[3] {
			fmt.Printf("Renderer Warning : Normal vectors not found (pinfo=%v)\n", pinfo)
		}
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
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, stride*4, offset*4)
		context.Call("enableVertexAttribArray", location)
		return nil
	}
	return fmt.Errorf("Invalid binding : attribute ('%s') failed to bind with buffer %T", dtype, buffer)
}
