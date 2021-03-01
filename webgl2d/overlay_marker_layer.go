package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

// ----------------------------------------------------------------------------
// OverlayMarkerLayer
// ----------------------------------------------------------------------------

type OverlayMarkerLayer struct {
	wctx    *common.WebGLContext //
	Markers []*SceneObject       // list of OverlayMarkers to be rendered (in pixels in CAMERA space)
}

func NewOverlayMarkerLayer(wctx *common.WebGLContext) *OverlayMarkerLayer {
	self := OverlayMarkerLayer{wctx: wctx}
	self.Markers = make([]*SceneObject, 0)
	return &self
}

func (self *OverlayMarkerLayer) Render(wctx *common.WebGLContext, pvm *geom2d.Matrix3) {
	renderer := NewRenderer(wctx)
	for _, marker := range self.Markers {
		renderer.RenderSceneObject(marker, pvm)
	}
}

// ----------------------------------------------------------------------------
// Arrow Markers
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) NewArrowMarker() *SceneObject {
	// Geometry (ARROW pointing left, with head at (0,0), and PIXEL length 10)
	geometry := NewGeometry()
	geometry.SetVertices([][2]float32{{0, 0}, {0.5, -0.3}, {0.5, -0.15}, {1, -0.15}, {1, 0.15}, {0.5, 0.15}, {0.5, 0.3}})
	geometry.SetFaces([][]uint32{{0, 1, 2, 3, 4, 5, 6}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 3, 4, 5, 6, 0}})
	geometry.Scale(30, 30).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	// Material (it's needed so that the same shader is shared for drawing EDGES and FACES)
	material := common.NewMaterial(self.wctx, "#aaffaa").SetColorForDrawMode(2, "#ff0000") // RED outline, PINK interior
	// Shader   (instance position 'pos' gives the origin, vertex XY coordinates 'oxy' gives offset in PIXEL)
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;		// Projection * View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 pos;		// world XY coordinates of the origin
		attribute vec2 oxy;		// vertex XY coordinates (pixels in CAMERA space)
		void main() {
			vec3 origin = pvm * vec3(pos, 1.0);
			vec2 offset = vec2(oxy.x * 2.0 / asp[0], oxy.y * 2.0 / asp[0]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := common.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")        // automatic binding of Proj*View*Model matrix
	shader.InitBindingForUniform("asp", "vec2", "renderer.aspect")     // automatic binding of material color
	shader.InitBindingForUniform("color", "vec3", "material.color")    // automatic binding of material color
	shader.InitBindingForAttribute("pos", "vec2", "instance.pose:2:0") // automatic binding of instance pose (:<stride>:<offset>)
	shader.InitBindingForAttribute("oxy", "vec2", "geometry.coords")   // automatic binding of point coordinates
	shader.CheckBindings()                                             // check validity of the shader
	marker := NewSceneObject(geometry, material, nil, shader, shader)  // shader shared for drawing EDGES & FACES
	return marker
}

func (self *OverlayMarkerLayer) AddArrowMarkersToTest() *OverlayMarkerLayer {
	arrow_poses := NewSceneObjectPoses(2, 3, []float32{0, 0, 40, 40, 40, 80}) // SceneObjectPoses
	arrow_marker := self.NewArrowMarker().SetInstancePoses(arrow_poses)       // SceneObject
	self.Markers = append(self.Markers, arrow_marker)
	return self
}

// ----------------------------------------------------------------------------
// Other Markers
// ----------------------------------------------------------------------------
