package webgl2d

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/wcommon/geom2d"
)

type Renderer struct {
	wctx *wcommon.WebGLContext
	axes *SceneObject
}

func NewRenderer(wctx *wcommon.WebGLContext) *Renderer {
	renderer := Renderer{wctx: wctx, axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *Renderer) Clear(scene *Scene) {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	rgb := scene.GetBkgColor()
	context.Call("clearColor", rgb[0], rgb[1], rgb[2], 1.0) // Set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT)       // clear the canvas
	context.Call("clear", constants.DEPTH_BUFFER_BIT)       // clear the canvas
}

func (self *Renderer) RenderAxes(camera *Camera, length float32) {
	if self.axes == nil {
		self.axes = NewSceneObject_2DAxes(self.wctx, length)
	}
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	context.Call("bindBuffer", constants.ARRAY_BUFFER, js.Null())
	context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, js.Null())
	self.RenderSceneObject(self.axes, &camera.pjvwmatrix) // (Proj * View) matrix
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(scene *Scene, camera *Camera) {
	// Render all the scene objects
	for _, sobj := range scene.objects {
		pvm_matrix := camera.pjvwmatrix.MultiplyToTheRight(&sobj.modelmatrix)
		self.RenderSceneObject(sobj, pvm_matrix) // (Proj * View * Model) matrix
	}
	// Render all the OverlayLayers
	for _, overlay := range scene.overlays {
		overlay.Render(&camera.pjvwmatrix)
	}
}

// ----------------------------------------------------------------------------
// Rendering SceneObject
// ----------------------------------------------------------------------------

func (self *Renderer) RenderSceneObject(sobj *SceneObject, pvm *geom2d.Matrix3) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	// Set DepthTest & Blending options
	if sobj.UseDepth {
		context.Call("enable", constants.DEPTH_TEST) // Enable depth test
		context.Call("depthFunc", constants.LEQUAL)  // Near things obscure far things
	} else {
		context.Call("disable", constants.DEPTH_TEST) // Disable depth test
	}
	if sobj.UseBlend {
		context.Call("enable", constants.BLEND)                                 // for pre-multiplied alpha
		context.Call("blendFunc", constants.ONE, constants.ONE_MINUS_SRC_ALPHA) // for pre-multiplied alpha
		// context.Call("blendFunc", constants.SRC_ALPHA, constants.ONE_MINUS_SRC_ALPHA) // for non pre-multiplied alpha
	} else {
		context.Call("disable", constants.BLEND) // Disable blending
	}
	// If necessary, then build WebGLBuffers for the SceneObject's Geometry
	if sobj.Geometry.IsDataBufferReady() == false {
		return errors.New("Failed to RenderSceneObject() : empty geometry data buffer")
	}
	if sobj.Geometry.IsWebGLBufferReady() == false {
		sobj.Geometry.BuildWebGLBuffers(self.wctx, true, true, true)
	}
	if sobj.poses != nil && sobj.poses.IsWebGLBufferReady() == false {
		sobj.poses.BuildWebGLBuffer(self.wctx)
		if !self.wctx.IsExtensionReady("ANGLE") {
			self.wctx.SetupExtension("ANGLE")
		}
	}
	// R3: Render the object with FACE shader
	if sobj.FShader != nil {
		err := self.render_scene_object_with_shader(sobj, pvm, 3, sobj.FShader)
		if err != nil {
			return err
		}
	}
	// R2: Render the object with EDGE shader
	if sobj.EShader != nil {
		err := self.render_scene_object_with_shader(sobj, pvm, 2, sobj.EShader)
		if err != nil {
			return err
		}
	}
	// R1: Render the object with VERTEX shader
	if sobj.VShader != nil {
		err := self.render_scene_object_with_shader(sobj, pvm, 1, sobj.VShader)
		if err != nil {
			return err
		}
	}
	// Render all the children
	for _, child := range sobj.children {
		new_pvm := pvm.MultiplyToTheRight(&child.modelmatrix)
		self.RenderSceneObject(child, new_pvm)
	}
	return nil
}

func (self *Renderer) render_scene_object_with_shader(sobj *SceneObject, pvm *geom2d.Matrix3, draw_mode int, shader *wcommon.Shader) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	// 1. Decide which Shader to use
	if shader == nil {
		return errors.New("Failed to RenderSceneObject() : shader not found")
	}
	context.Call("useProgram", shader.GetShaderProgram())
	// 2. bind the uniforms of the shader program
	for uname, umap := range shader.GetUniformBindings() {
		if err := self.bind_uniform(uname, umap, draw_mode, sobj.Material, pvm); err != nil {
			if err.Error() != "Texture is not ready" {
				fmt.Println(err.Error())
			}
			return err
		}
	}
	// 3. bind the attributes of the shader program
	for aname, amap := range shader.GetAttributeBindings() {
		if err := self.bind_attribute(aname, amap, draw_mode, sobj.Geometry, sobj.poses); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 4. draw  (Note that ARRAY_BUFFER was binded already in the attribut-binding step)
	switch draw_mode {
	case 3: // draw TRIANGLES (FACES)
		buffer, count, _ := sobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
			if sobj.poses == nil {
				context.Call("drawElements", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.Count
				ext.Call("drawElementsInstancedANGLE", constants.TRIANGLES, count, constants.UNSIGNED_INT, 0, pose_count)
			}
		}
	case 2: // draw LINES (EDGES)
		buffer, count, _ := sobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, buffer)
			if sobj.poses == nil {
				context.Call("drawElements", constants.LINES, count, constants.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.Count
				ext.Call("drawElementsInstancedANGLE", constants.LINES, count, constants.UNSIGNED_INT, 0, pose_count)
			}
		}
	case 1: // draw POINTS (VERTICES)
		_, count, pinfo := sobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			vert_count := count / pinfo[0] // number of vertices
			if sobj.poses == nil {
				context.Call("drawArrays", constants.POINTS, 0, vert_count) // (mode, first, count)
			} else {
				ext, pose_count := self.wctx.GetExtension("ANGLE"), sobj.poses.Count
				ext.Call("drawArraysInstancedANGLE", constants.POINTS, 0, vert_count, pose_count)
			}
		}
	default:
		err := fmt.Errorf("Unknown mode to draw : %s\n", draw_mode)
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func (self *Renderer) bind_uniform(uname string, umap map[string]interface{},
	draw_mode int, material *wcommon.Material, pvm *geom2d.Matrix3) error {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	if umap["location"] == nil {
		err := errors.New("Failed to bind uniform : call 'shader.CheckBinding()' before rendering")
		return err
	}
	location, dtype := umap["location"].(js.Value), umap["dtype"].(string)
	if umap["autobinding"] != nil {
		autobinding := umap["autobinding"].(string)
		autobinding_split := strings.Split(autobinding, ":")
		autobinding0 := autobinding_split[0]
		switch autobinding0 {
		case "material.color":
			c := [4]float32{0, 1, 1, 1}
			if material != nil {
				c = material.GetDrawModeColor(draw_mode) // get color from material (for the DrawMode)
			}
			switch dtype {
			case "vec3":
				context.Call("uniform3f", location, c[0], c[1], c[2])
				return nil
			case "vec4":
				context.Call("uniform4f", location, c[0], c[1], c[2], c[3])
				return nil
			}
		case "material.texture":
			if material == nil || !material.IsTextureReady() || material.IsTextureLoading() {
				return errors.New("Texture is not ready")
			}
			txt_unit := 0
			if len(autobinding_split) >= 2 {
				txt_unit, _ = strconv.Atoi(autobinding_split[1])
			}
			texture_unit := js.ValueOf(constants.TEXTURE0.Int() + txt_unit)
			context.Call("activeTexture", texture_unit)                              // activate texture unit N
			context.Call("bindTexture", constants.TEXTURE_2D, material.GetTexture()) // bind the texture
			context.Call("uniform1i", location, txt_unit)                            // give shader the unit number
			return nil
		case "renderer.aspect": // vec2
			wh := self.wctx.GetWH()
			context.Call("uniform2f", location, float32(wh[0]), float32(wh[1]))
			return nil
		case "renderer.pvm": // mat3
			elements := pvm.GetElements()
			e := wcommon.ConvertGoSliceToJsTypedArray(elements[:]) // ModelView matrix, converted to JavaScript 'Float32Array'
			context.Call("uniformMatrix3fv", location, false, e)   // gl.uniformMatrix3fv(location, transpose, values_array)
			return nil
		}
		return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, autobinding)
	} else if umap["value"] != nil {
		v := umap["value"].([]float32)
		switch dtype {
		case "int":
			context.Call("uniform1i", location, int(v[0]))
			return nil
		case "float":
			context.Call("uniform1f", location, v[0])
			return nil
		case "vec2":
			context.Call("uniform2f", location, v[0], v[1])
			return nil
		case "vec3":
			context.Call("uniform3f", location, v[0], v[1], v[2])
			return nil
		case "vec4":
			context.Call("uniform4f", location, v[0], v[1], v[2], v[3])
			return nil
		}
		return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, v)
	} else {
		return fmt.Errorf("Failed to bind uniform '%s' (%s)", uname, dtype)
	}
}

func (self *Renderer) bind_attribute(aname string, amap map[string]interface{},
	draw_mode int, geometry *Geometry, poses *wcommon.SceneObjectPoses) error {
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
	autobinding0 := autobinding_split[0]
	switch autobinding0 {
	case "geometry.coords": // 2 * float32 in 8 bytes (2 float32)
		buffer, _, pinfo := geometry.GetWebGLBuffer(1) // pinfo : [3]{stride, xy_offset, uv_offset}
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 2, constants.FLOAT, false, pinfo[0]*4, pinfo[1]*4)
		context.Call("enableVertexAttribArray", location)
		if self.wctx.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
		}
		return nil
	case "geometry.textuv": // 2 * uint16 in 4 bytes (1 float32)
		buffer, _, pinfo := geometry.GetWebGLBuffer(1) // pinfo : [3]{stride, xy_offset, uv_offset}
		context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer)
		context.Call("vertexAttribPointer", location, 2, constants.UNSIGNED_SHORT, true, pinfo[0]*4, pinfo[2]*4)
		context.Call("enableVertexAttribArray", location)
		if self.wctx.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
		}
		return nil
	case "instance.pose":
		if poses != nil && len(autobinding_split) == 3 { // it's like "instance.pose:<stride>:<offset>"
			size := get_count_from_type(dtype)
			stride, _ := strconv.Atoi(autobinding_split[1])
			offset, _ := strconv.Atoi(autobinding_split[2])
			context.Call("bindBuffer", constants.ARRAY_BUFFER, poses.WebGLBuffer)
			context.Call("vertexAttribPointer", location, size, constants.FLOAT, false, stride*4, offset*4)
			context.Call("enableVertexAttribArray", location)
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 1) // divisor == 1
			return nil
		}
	default:
		buffer, stride_i, offset_i := amap["buffer"], amap["stride"], amap["offset"]
		if buffer != nil && stride_i != nil && offset_i != nil {
			size, stride, offset := get_count_from_type(dtype), stride_i.(int), offset_i.(int)
			context.Call("bindBuffer", constants.ARRAY_BUFFER, buffer.(js.Value))
			context.Call("vertexAttribPointer", location, size, constants.FLOAT, false, stride*4, offset*4)
			context.Call("enableVertexAttribArray", location)
			if self.wctx.IsExtensionReady("ANGLE") {
				// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
				self.wctx.GetExtension("ANGLE").Call("vertexAttribDivisorANGLE", location, 0) // divisor == 0
			}
		}
	}
	return fmt.Errorf("Failed to bind attribute '%s' (%s) with %v", aname, dtype, amap)
}

func get_count_from_type(dtype string) int {
	switch dtype {
	case "float":
		return 1
	case "vec2":
		return 2
	case "vec3":
		return 3
	case "vec4":
		return 4
	default:
		return 0
	}
}
