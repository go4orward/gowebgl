package webgl3d

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	if sobj.poses != nil && sobj.poses.IsWebGLBufferReady() == false {
		sobj.poses.build_webgl_buffers(self.wctx)
		if !self.wctx.IsExtensionReady("ANGLE") {
			self.wctx.SetupExtension("ANGLE")
		}
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
		if err := self.bind_uniform(uname, umap, sobj.material, proj, viewm); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 4. bind the attributes of the shader program
	for aname, amap := range shader.GetAttributeBindings() {
		if err := self.bind_attribute(aname, amap, sobj.geometry, sobj.poses); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 5. draw
	for _, draw_mode := range shader.GetThingsToDraw() {
		// Note that ARRAY_BUFFER was binded already in the previous step (during attribute binding)
		switch draw_mode {
		case "POINTS", "VERTICES":
			_, count, _ := sobj.geometry.GetWebGLBuffer("POINTS")
			if count > 0 {
				if sobj.poses == nil {
					context.Call("drawArrays", constants.POINTS, 0, count) // (mode, first, count)
				} else {
					ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.GetPoseCount()
					ext.Call("drawArraysInstancedANGLE", constants.POINTS, 0, count, pose_count)
				}
			}
		case "LINES", "EDGES":
			buffer, count, _ := sobj.geometry.GetWebGLBuffer("LINES")
			if count > 0 {
				context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
				if sobj.poses == nil {
					context.Call("drawElements", constants.LINES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
				} else {
					ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.GetPoseCount()
					ext.Call("drawElementsInstancedANGLE", constants.LINES, count, constants.UNSIGNED_INT, 0, pose_count)
				}
			}
		case "TRIANGLES", "FACES":
			buffer, count, _ := sobj.geometry.GetWebGLBuffer("TRIANGLES")
			if count > 0 {
				context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
				if sobj.poses == nil {
					context.Call("drawElements", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
				} else {
					ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.GetPoseCount()
					ext.Call("drawElementsInstancedANGLE", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0, pose_count)
				}
			}
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

func (self *Renderer) bind_uniform(uname string, umap map[string]interface{}, material *Material, proj *geom3d.Matrix4, viewm *geom3d.Matrix4) error {
	context := self.wctx.GetContext()
	if umap["location"] == nil {
		err := errors.New("Failed to bind uniform : call 'shader.CheckBinding()' before rendering")
		return err
	}
	location, dtype := umap["location"].(js.Value), umap["dtype"].(string)
	autobinding := umap["autobinding"].(string)
	// fmt.Printf("Uniform (%s) : autobinding= '%s'\n", dtype, autobinding)
	switch autobinding {
	case "renderer.proj": // mat4
		e := (*proj.GetElements())[:]
		m := common.ConvertGoSliceToJsTypedArray(e)          // Projection matrix, converted to JavaScript 'Float32Array'
		context.Call("uniformMatrix4fv", location, false, m) // gl.uniformMatrix4fv(location, transpose, values_array)
		return nil
	case "renderer.vwmd": // mat4
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
		if material != nil {
			c = material.color
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
	default:
		value := umap["value"]
		if value != nil {
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
		}
	}
	return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, autobinding, umap)
}

func (self *Renderer) bind_attribute(aname string, amap map[string]interface{}, geometry *Geometry, poses *SceneObjectPoses) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	if amap["location"] == nil {
		err := errors.New("Failed to bind attribute : call 'shader.CheckBinding()' before rendering")
		return err
	}
	location, dtype := amap["location"].(js.Value), amap["dtype"].(string)
	autobinding := amap["autobinding"].(string)
	// fmt.Printf("Attribute (%s) : autobinding= '%s'\n", dtype, autobinding)
	autobinding_split := strings.Split(autobinding, ":")
	autobinding = autobinding_split[0]
	switch autobinding {
	case "geometry.coords":
		buffer, _, pinfo := geometry.GetWebGLBuffer("POINTS")
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 3, constants.FLOAT, false, pinfo[0]*4, pinfo[1]*4)
		context.Call("enableVertexAttribArray", location)
		if poses != nil { // context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
		}
		return nil
	case "geometry.textuv":
		buffer, _, pinfo := geometry.GetWebGLBuffer("POINTS")
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, pinfo[0]*4, pinfo[2]*4)
		context.Call("enableVertexAttribArray", location)
		if pinfo[1] == pinfo[2] {
			fmt.Printf("Renderer Warning : Texture UV coordinates not found (pinfo=%v)\n", pinfo)
		}
		if poses != nil { // context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
		}
		return nil
	case "geometry.normal":
		buffer, _, pinfo := geometry.GetWebGLBuffer("POINTS")
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 3, constants.FLOAT, false, pinfo[0]*4, pinfo[3]*4)
		context.Call("enableVertexAttribArray", location)
		if pinfo[1] == pinfo[3] {
			fmt.Printf("Renderer Warning : Normal vectors not found (pinfo=%v)\n", pinfo)
		}
		if poses != nil { // context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
		}
		return nil
	case "instance.pose":
		if poses != nil && len(autobinding_split) == 3 { // it's like "instance.pose:<stride>:<offset>"
			count := get_count_from_type(dtype)
			stride, _ := strconv.Atoi(autobinding_split[1])
			offset, _ := strconv.Atoi(autobinding_split[2])
			context.Call("bindBuffer", constants.ARRAY_BUFFER, poses.webgl_buffer)
			context.Call("vertexAttribPointer", location, count, constants.FLOAT, false, stride*4, offset*4)
			context.Call("enableVertexAttribArray", location)
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 1) // divisor == 1
			return nil
		}
	default:
		buffer, stride_i, offset_i := amap["buffer"], amap["stride"], amap["offset"]
		if buffer != nil && stride_i != nil && offset_i != nil {
			count, stride, offset := get_count_from_type(dtype), stride_i.(int), offset_i.(int)
			context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer.(js.Value))
			context.Call("vertexAttribPointer", location, count, constants.FLOAT, false, stride*4, offset*4)
			context.Call("enableVertexAttribArray", location)
			if poses != nil { // context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
				self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
			}
		}
	}
	return fmt.Errorf("Failed to bind attribute '%s' (%s) with %v", aname, dtype, amap)
}

func get_count_from_type(dtype string) int {
	switch dtype {
	case "vec2":
		return 2
	case "vec3":
		return 3
	case "vec4":
		return 4
	case "mat2":
		return 4
	case "mat3":
		return 9
	case "mat4":
		return 16
	default:
		return 0
	}
}
