package webgl3d

import (
	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/wcommon/geom3d"
	"github.com/go4orward/gowebgl/webgl2d"
)

type OverlayMarkerLayer struct {
	wctx    *wcommon.WebGLContext //
	Markers []*SceneObject        // list of OverlayMarkers to be rendered (in pixels in CAMERA space)
}

func NewOverlayMarkerLayer(wctx *wcommon.WebGLContext) *OverlayMarkerLayer {
	self := OverlayMarkerLayer{wctx: wctx}
	self.Markers = make([]*SceneObject, 0)
	return &self
}

func (self *OverlayMarkerLayer) Render(proj *geom3d.Matrix4, vwmd *geom3d.Matrix4) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.wctx)
	for _, marker := range self.Markers {
		if marker.poses != nil {
			renderer.RenderSceneObject(marker, proj, vwmd)
		} else {
			vwmd := vwmd.MultiplyToTheRight(&marker.modelmatrix)
			renderer.RenderSceneObject(marker, proj, vwmd)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Markers
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) AddMarker(marker ...*SceneObject) *OverlayMarkerLayer {
	for i := 0; i < len(marker); i++ {
		self.Markers = append(self.Markers, marker[i])
	}
	return self
}

func (self *OverlayMarkerLayer) AddArrowMarker(size float32, color string, outline_color string, rotation float32, xyz [3]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : arrow := layer.CreateMarkerArrow();  layer.AddLabel(label)
	arrow := self.CreateArrowMarker(size, color, outline_color, false)
	arrow.Rotate([3]float32{0, 0, 1}, rotation).Translate(xyz[0], xyz[1], xyz[2])
	self.Markers = append(self.Markers, arrow)
	return self
}

func (self *OverlayMarkerLayer) AddArrowHeadMarker(size float32, color string, outline_color string, rotation float32, xyz [3]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : ahead := layer.CreateMarkerArrowHead();  ahead.Translate();  layer.AddLabel(ahead)
	ahead := self.CreateArrowHeadMarker(size, color, outline_color, false)
	ahead.Rotate([3]float32{0, 0, 1}, rotation).Translate(xyz[0], xyz[1], xyz[2])
	self.Markers = append(self.Markers, ahead)
	return self
}

func (self *OverlayMarkerLayer) AddArrowHeadMarkerForTest() *OverlayMarkerLayer {
	if true {
		ahead1 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", false)
		ahead2 := self.CreateArrowHeadMarker(20, "#aaffaa", "#00ff00", false).Translate(0.5, 0.5, 0.5)
		return self.AddMarker(ahead1, ahead2)
	} else {
		ahead := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", true)
		poses := wcommon.NewSceneObjectPoses(3, 3, []float32{0, 0, -1, 1, 1, 1, 0, 0, 1})
		return self.AddMarker(ahead.SetInstancePoses(poses))
	}
}

// ----------------------------------------------------------------------------
// Marker Examples
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateArrowMarker(size float32, color string, outline_color string, use_poses bool) *SceneObject {
	geometry := webgl2d.NewGeometry() // ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {0.5, -0.3}, {0.5, -0.15}, {1, -0.15}, {1, 0.15}, {0.5, 0.15}, {0.5, 0.3}})
	geometry.SetFaces([][]uint32{{0, 1, 2, 3, 4, 5, 6}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 3, 4, 5, 6, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material := wcommon.NewMaterial(self.wctx, color).SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(use_poses)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

func (self *OverlayMarkerLayer) CreateArrowHeadMarker(size float32, color string, outline_color string, use_poses bool) *SceneObject {
	geometry := webgl2d.NewGeometry() // ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {1, -0.6}, {1, +0.6}})
	geometry.SetFaces([][]uint32{{0, 1, 2}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material := wcommon.NewMaterial(self.wctx, color).SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(use_poses)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

// ----------------------------------------------------------------------------
// Shader for Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) GetShaderForMarker(use_poses bool) *wcommon.Shader {
	var shader *wcommon.Shader = nil
	if !use_poses { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat4 proj;	// 3D Projection matrix
		uniform   mat4 vwmd;	// 3D View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 gvxy;	// 2D vertex XY coordinates (offset; pixels in CAMERA space)
		void main() {
			vec4 origin = proj * vwmd * vec4(0.0, 0.0, 0.0, 1.0);
			origin = origin / origin.w;
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")     // Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")     // View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")    // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")   // material color
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords") // 2D offset coordinates (in CAMERA space)
	} else { // Shader for multiple instance poses ('orgn')
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat4 proj;	// 3D Projection matrix
		uniform   mat4 vwmd;	// 3D View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec3 orgn;	// 3D world XYZ coordinates of the origin
		attribute vec2 gvxy;	// 2D vertex XY coordinates (offset; pixels in CAMERA space)
		void main() {
			vec4 origin = proj * vwmd * vec4(orgn, 1.0);
			origin = origin / origin.w;
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")       // 3D Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")       // 3D View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")      // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")     // material color
		shader.SetBindingForAttribute("orgn", "vec3", "instance.pose:3:0") // instance pose (:<stride>:<offset>)
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")   // 2D offset coordinates (in CAMERA space)
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}
