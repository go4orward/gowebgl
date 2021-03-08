package webgl2d

import (
	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/wcommon/geom2d"
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

func (self *OverlayMarkerLayer) Render(pvm *geom2d.Matrix3) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.wctx)
	for _, marker := range self.Markers {
		if marker.poses != nil {
			renderer.RenderSceneObject(marker, pvm)
		} else {
			new_pvm := pvm.MultiplyToTheRight(&marker.modelmatrix)
			renderer.RenderSceneObject(marker, new_pvm)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Markers
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) AddMarker(marker *SceneObject) *OverlayMarkerLayer {
	self.Markers = append(self.Markers, marker)
	return self
}

func (self *OverlayMarkerLayer) AddMarkerArrow(size float32, color string, outline_color string, rotation float32, xy [2]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : arrow := layer.CreateMarkerArrow();  layer.AddLabel(label)
	arrow := self.CreateMarkerArrow(size, color, outline_color, false)
	arrow.Rotate(rotation).Translate(xy[0], xy[1])
	self.Markers = append(self.Markers, arrow)
	return self
}

func (self *OverlayMarkerLayer) AddMarkerArrowHead(size float32, color string, outline_color string, rotation float32, xy [2]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : ahead := layer.CreateMarkerArrowHead();  ahead.Translate();  layer.AddLabel(ahead)
	ahead := self.CreateMarkerArrowHead(size, color, outline_color, false)
	ahead.Rotate(rotation).Translate(xy[0], xy[1])
	self.Markers = append(self.Markers, ahead)
	return self
}

func (self *OverlayMarkerLayer) AddMarkerArrowHeadsToTest() *OverlayMarkerLayer {
	ahead := self.CreateMarkerArrowHead(20, "#ffaaaa", "#ff0000", true)
	poses := wcommon.NewSceneObjectPoses(2, 3, []float32{0, 0, 40, 40, 40, 80})
	return self.AddMarker(ahead.SetInstancePoses(poses))
}

// ----------------------------------------------------------------------------
// Marker Examples
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateMarkerArrow(size float32, color string, outline_color string, multiple bool) *SceneObject {
	geometry := NewGeometry() // ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {0.5, -0.3}, {0.5, -0.15}, {1, -0.15}, {1, 0.15}, {0.5, 0.15}, {0.5, 0.3}})
	geometry.SetFaces([][]uint32{{0, 1, 2, 3, 4, 5, 6}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 3, 4, 5, 6, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material := wcommon.NewMaterial(self.wctx, color).SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(multiple)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

func (self *OverlayMarkerLayer) CreateMarkerArrowHead(size float32, color string, outline_color string, multiple bool) *SceneObject {
	geometry := NewGeometry() // ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {1, -0.6}, {1, +0.6}})
	geometry.SetFaces([][]uint32{{0, 1, 2}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material := wcommon.NewMaterial(self.wctx, color).SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(multiple)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

// ----------------------------------------------------------------------------
// Shader for Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) GetShaderForMarker(multiple bool) *wcommon.Shader {
	var shader *wcommon.Shader = nil
	if !multiple { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;		// Projection * View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 gvxy;	// vertex XY coordinates (pixels in CAMERA space)
		void main() {
			vec3 origin = pvm * vec3(0.0, 0.0, 1.0);
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[0]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")       // Proj*View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")    // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")   // material color
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords") // offset coordinates (in CAMERA space)
	} else { // Shader for multiple instance poses ('orgn')
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;		// Projection * View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 orgn;	// world XY coordinates of the origin
		attribute vec2 gvxy;	// vertex XY coordinates (pixels in CAMERA space)
		void main() {
			vec3 origin = pvm * vec3(orgn, 1.0);
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[0]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")         // automatic binding of Proj*View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")      // automatic binding of AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")     // automatic binding of material color
		shader.SetBindingForAttribute("orgn", "vec2", "instance.pose:2:0") // automatic binding of instance pose (:<stride>:<offset>)
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")   // automatic binding of point coordinates
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}
