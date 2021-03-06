package webgl2d

import (
	"github.com/go4orward/gowebgl/geom2d"
	"github.com/go4orward/gowebgl/wcommon"
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

func (self *OverlayMarkerLayer) AddMarker(marker ...*SceneObject) *OverlayMarkerLayer {
	for i := 0; i < len(marker); i++ {
		self.Markers = append(self.Markers, marker[i])
	}
	return self
}

func (self *OverlayMarkerLayer) AddArrowMarker(size float32, color string, outline_color string, rotation float32, xy [2]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a ARROW marker,
	//   which is equivalent to : arrow := layer.CreateArrowMarker();  layer.AddMarker(label)
	arrow := self.CreateArrowMarker(size, color, outline_color, false)
	arrow.Rotate(rotation).Translate(xy[0], xy[1])
	self.AddMarker(arrow)
	return self
}

func (self *OverlayMarkerLayer) AddArrowHeadMarker(size float32, color string, outline_color string, rotation float32, xy [2]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a ARROW_HEAD marker,
	//   which is equivalent to : ahead := layer.CreateArrowHeadMarker();  ahead.Translate();  layer.AddMarker(ahead)
	ahead := self.CreateArrowHeadMarker(size, color, outline_color, false)
	ahead.Rotate(rotation).Translate(xy[0], xy[1])
	self.AddMarker(ahead)
	return self
}

func (self *OverlayMarkerLayer) AddSpriteMarker(imgpath string, color string, wh [2]float32, xy [2]float32, offref string) *OverlayMarkerLayer {
	// Convenience function to quickly add a SPRITE marker,
	//   which is equivalent to : sprite := layer.CreateSpriteMarker();  sprite.Translate();  layer.AddMarker(sprite)
	sprite := self.CreateSpriteMarker(imgpath, color, wh, offref, false).Translate(xy[0], xy[1])
	return self.AddMarker(sprite)
}

func (self *OverlayMarkerLayer) AddMarkersForTest() *OverlayMarkerLayer {
	ahead1 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", false)
	ahead2 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", false).Translate(40, 80)
	ahead3 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", true)
	ahead3.SetupPoses(2, 4, []float32{20, 90, 30, 90, 40, 90, 50, 90})
	sprite := self.CreateSpriteMarker("/assets/map_marker.png", "#ff0000", [2]float32{20, 20}, "M_BTM", false)
	return self.AddMarker(ahead1, ahead2, ahead3, sprite)
}

// ----------------------------------------------------------------------------
// Arrow Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateArrowMarker(size float32, color string, outline_color string, use_poses bool) *SceneObject {
	geometry := NewGeometry() // ARROW pointing left, with tip at (0,0)
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
	geometry := NewGeometry() // ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {1, -0.6}, {1, +0.6}})
	geometry.SetFaces([][]uint32{{0, 1, 2}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material := wcommon.NewMaterial(self.wctx, color).SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(use_poses)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

func (self *OverlayMarkerLayer) GetShaderForMarker(use_poses bool) *wcommon.Shader {
	var shader *wcommon.Shader = nil
	if !use_poses { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;		// Projection * View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 gvxy;	// vertex XY coordinates (pixels in CAMERA space)
		void main() {
			vec3 origin = pvm * vec3(0.0, 0.0, 1.0);
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
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
	} else { // Shader for multiple instance poses ('iorg')
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;		// Projection * View * Model matrix
		uniform   vec2 asp;		// aspect ratio, w : h
		attribute vec2 iorg;	// world XY coordinates of the origin
		attribute vec2 gvxy;	// vertex XY coordinates (pixels in CAMERA space)
		void main() {
			vec3 origin = pvm * vec3(iorg, 1.0);
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")         // Proj*View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")      // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")     // material color
		shader.SetBindingForAttribute("iorg", "vec2", "instance.pose:2:0") // instance pose (:<stride>:<offset>)
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")   // point coordinates
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}

// ----------------------------------------------------------------------------
// Sprite Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateSpriteMarker(imgpath string, color string, wh [2]float32, offref string, use_poses bool) *SceneObject {
	geometry := NewGeometry_Origin() // geometry with only one vertex at (0,0)
	material := wcommon.NewMaterial(self.wctx, imgpath).SetColorForDrawMode(0, color)
	// wh := [2]float32{float32(material.GetTextureWH()[0]), float32(material.GetTextureWH()[1])}
	var offrot [3]float32
	switch offref {
	case "L_TOP":
		offrot = [3]float32{+wh[0] / 2, -wh[1] / 2, 0}
	case "M_TOP":
		offrot = [3]float32{0, -wh[1] / 2, 0}
	case "R_TOP":
		offrot = [3]float32{-wh[0] / 2, -wh[1] / 2, 0}
	case "L_CTR":
		offrot = [3]float32{+wh[0] / 2, 0, 0}
	case "M_CTR", "CENTER":
		offrot = [3]float32{0, 0, 0}
	case "R_CTR":
		offrot = [3]float32{-wh[0] / 2, 0, 0}
	case "L_BTM":
		offrot = [3]float32{+wh[0] / 2, +wh[1] / 2, 0}
	case "M_BTM":
		offrot = [3]float32{0, +wh[1] / 2, 0}
	case "R_BTM":
		offrot = [3]float32{-wh[0] / 2, +wh[1] / 2, 0}
	default:
	}
	shader := self.GetShaderForSpriteMarker(wh, offrot, use_poses)
	sprite := NewSceneObject(geometry, material, shader, nil, nil)
	sprite.UseBlend = true
	return sprite
}

func (self *OverlayMarkerLayer) GetShaderForSpriteMarker(wh [2]float32, offrot [3]float32, use_poses bool) *wcommon.Shader {
	var shader *wcommon.Shader = nil
	if !use_poses { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat3  pvm;		// Projection * View * Model matrix
			uniform   vec2  asp;		// aspect ratio, w : h
			uniform   vec2  wh;			// size of the sprite
			uniform   vec3  offr;		// offset of the label (CAMERA XY in pixel) & rotation_angle
			attribute vec2  gvxy;		// vertex XY coordinates (pixels in CAMERA space)
			void main() {
				vec3 origin = pvm * vec3(0.0, 0.0, 1.0);
				vec2 offset = vec2((offr.x + gvxy.x) * 2.0 / asp[0], (offr.y + gvxy.y) * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
				if (wh[0] > wh[1]) {
					gl_PointSize = wh[0];	// sprite size is its width
				} else {
					gl_PointSize = wh[1];	// sprite size is its height
				}
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
			uniform   vec4  color;		// color of the sprite
			uniform   vec2  wh;			// size of the sprite
			void main() {
				vec2 uv = gl_PointCoord;
				if (wh[0] > wh[1]) {
					uv[1] = (uv[1] - 0.5) * wh[0]/wh[1] + 0.5;
				} else {
					uv[0] = (uv[0] - 0.5) * wh[1]/wh[0] + 0.5;
				}
				if (uv[0] < 0.0 || uv[0] > 1.0) discard;
				if (uv[1] < 0.0 || uv[1] > 1.0) discard;
				gl_FragColor = texture2D(texture, uv) * color;
			}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")              // Proj*View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
		shader.SetBindingForUniform("wh", "vec2", wh[:])                        // sprite size
		shader.SetBindingForUniform("offr", "vec3", offrot[:])                  // sprite offset & rotation angle
		shader.SetBindingForUniform("color", "vec4", "material.color")          // color to be multiplied with sprite texture
		shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")        // offset coordinates (in CAMERA space)
	} else { // Shader for multiple instance poses ('ixy')
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat3  pvm;		// Projection * View * Model matrix
			uniform   vec2  asp;		// aspect ratio, w : h
			uniform   vec2  wh;			// size of the sprite
			uniform   vec3  offr;		// offset of the label (CAMERA XY in pixel) & rotation_angle
			attribute vec2  gvxy;		// vertex XY coordinates (pixels in CAMERA space)
			attribute vec2  ixy;		// sprite instance position (in WORLD XY)
			void main() {
				vec3 origin = pvm * vec3(ixy.x, ixy.y, 1.0);
				vec2 offset = vec2((offr.x + gvxy.x) * 2.0 / asp[0], (offr.y + gvxy.y) * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, 0.0, 1.0);
				if (wh[0] > wh[1]) {
					gl_PointSize = wh[0];	// sprite size is its width
				} else {
					gl_PointSize = wh[1];	// sprite size is its height
				}
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
			uniform   vec4  color;		// color of the sprite
			uniform   vec2  wh;			// size of the sprite
			void main() {
				vec2 uv = gl_PointCoord;
				if (wh[0] > wh[1]) {
					uv[1] = (uv[1] - 0.5) * wh[0]/wh[1] + wh[1]/2.0;
				} else {
					uv[0] = (uv[0] - 0.5) * wh[1]/wh[0] + wh[0]/2.0;
				}
				if (uv[0] < 0.0 || uv[0] > 1.0) discard;
				if (uv[1] < 0.0 || uv[1] > 1.0) discard;
				gl_FragColor = texture2D(texture, uv) * color;
			}`
		shader, _ = wcommon.NewShader(self.wctx, vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")              // Proj*View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
		shader.SetBindingForUniform("wh", "vec2", wh[:])                        // sprite size
		shader.SetBindingForUniform("offr", "vec3", offrot[:])                  // sprite offset & rotation angle
		shader.SetBindingForUniform("color", "vec4", "material.color")          // color to be multiplied with sprite texture
		shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
		shader.SetBindingForAttribute("ixy", "vec2", "instance.pose:2:0")       // sprite instance position (in WORLD XY)
		shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")        // offset coordinates (in CAMERA space)
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}
