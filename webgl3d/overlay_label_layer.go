package webgl3d

import (
	"fmt"
	"strings"

	"github.com/go4orward/gowebgl/geom3d"
	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/webgl2d"
)

type OverlayLabel struct {
	text    string       // text of the label
	xyz     [3]float32   // origin of the label (in WORLD space)
	chwh    [2]float32   // character width & height
	color   string       // color of the label (like "#ff0000")
	offset  [2]float32   // offset from origin (in pixels in CAMERA space)
	offref  string       // offset reference type, like "L_TOP", "R_BTM", "CENTER", etc
	angle   float32      // rotation angle
	bkgtype string       // background type, like "box:#ffff00:#000000", "under:#000000", etc
	txtobj  *SceneObject // SceneObject for rendering the label text
	bkgobj  *SceneObject // SceneObject for rendering the background
}

type OverlayLabelLayer struct {
	wctx             *wcommon.WebGLContext //
	alphabet_texture *wcommon.Material     // Alphabet texture
	Labels           []*OverlayLabel       //
}

func NewOverlayLabelLayer(wctx *wcommon.WebGLContext, fontsize int, outlined bool) *OverlayLabelLayer {
	self := OverlayLabelLayer{wctx: wctx} // let 'fontsize' of ALPHABET texture to be 20, by default
	self.alphabet_texture = wcommon.NewMaterial_AlphabetTexture(wctx, fontsize, "#ffffff", outlined)
	self.Labels = make([]*OverlayLabel, 0)
	return &self
}

func (self *OverlayLabelLayer) Render(proj *geom3d.Matrix4, view *geom3d.Matrix4) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.wctx)
	for _, label := range self.Labels {
		if label.bkgobj != nil {
			renderer.RenderSceneObject(label.bkgobj, proj, view)
		}
		if label.txtobj != nil {
			renderer.RenderSceneObject(label.txtobj, proj, view)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Labels
// ----------------------------------------------------------------------------

func (self *OverlayLabelLayer) AddLabel(labels ...*OverlayLabel) *OverlayLabelLayer {
	for i := 0; i < len(labels); i++ {
		label := labels[i]
		if label.txtobj == nil && label.text != "" {
			label.build_labeltext_object(self.wctx, self.alphabet_texture)
		}
		if label.bkgobj == nil && label.bkgtype != "" {
			label.build_background_object(self.wctx)
		}
		self.Labels = append(self.Labels, label)
	}
	return self
}

func (self *OverlayLabelLayer) FindLabel(label_text string) *OverlayLabel {
	for _, label := range self.Labels {
		if label.text == label_text {
			return label
		}
	}
	return nil
}

func (self *OverlayLabelLayer) CreateLabel(label_text string, xyz [3]float32, color string) *OverlayLabel {
	chwh := self.alphabet_texture.GetAlaphabetCharacterWH(1.0)
	label := &OverlayLabel{text: label_text, xyz: xyz, chwh: chwh, color: color}
	return label
}

func (self *OverlayLabelLayer) AddTextLabel(label_text string, xyz [3]float32, color string, offref string) *OverlayLabelLayer {
	// Convenience function to quickly add a Label,
	//   which simplifies all the following steps:
	//   label := layer.CreateLabel();  label.SetPose();  layer.AddLabel(label)
	chwh := self.alphabet_texture.GetAlaphabetCharacterWH(1.0)
	label := &OverlayLabel{text: label_text, xyz: xyz, chwh: chwh, color: color}
	label.SetPose(0, offref, [2]float32{0, 0})
	self.AddLabel(label)
	return self
}

func (self *OverlayLabelLayer) AddLabelsForTest() *OverlayLabelLayer {
	self.AddTextLabel("AhjgyZ", [3]float32{0, .5, .5}, "#ff0000", "")
	label2 := self.CreateLabel("Hello!", [3]float32{0, 1, 1}, "#0000ff")
	label2.SetPose(0, "L_BTM", [2]float32{30, 30}).SetBackground("under:#000000")
	return self.AddLabel(label2)
}

// ----------------------------------------------------------------------------
// Label Functions
// ----------------------------------------------------------------------------

func (self *OverlayLabel) SetCharacterWH(chwh [2]float32) *OverlayLabel {
	self.chwh = chwh // character width & height
	return self
}

func (self *OverlayLabel) SetPose(rotation float32, offset_reference string, offset [2]float32) *OverlayLabel {
	text_width := self.chwh[0] * float32(len([]rune(self.text)))
	text_height := self.chwh[1]
	self.offref = offset_reference
	switch self.offref {
	case "L_TOP":
		self.offset = [2]float32{offset[0], offset[1] - text_height/2}
	case "L_MID":
		self.offset = [2]float32{offset[0], offset[1]}
	case "L_BTM":
		self.offset = [2]float32{offset[0], offset[1] + text_height/2}
	case "M_TOP":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1] - text_height/2}
	case "CENTER":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1]}
	case "M_BTM":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1] + text_height/2}
	case "R_TOP":
		self.offset = [2]float32{offset[0] - text_width, offset[1] - text_height/2}
	case "R_MID":
		self.offset = [2]float32{offset[0] - text_width, offset[1]}
	case "R_BTM":
		self.offset = [2]float32{offset[0] - text_width, offset[1] + text_height/2}
	default: // same as "L_MID"
		self.offset = [2]float32{offset[0], offset[1]}
	}
	return self
}

func (self *OverlayLabel) SetBackground(bkgtype string) *OverlayLabel {
	self.bkgtype = bkgtype
	return self
}

// ----------------------------------------------------------------------------
// Building SceneObjects
// ----------------------------------------------------------------------------

func (self *OverlayLabel) build_labeltext_object(wctx *wcommon.WebGLContext, alphabet_texture *wcommon.Material) *OverlayLabel {
	geometry := webgl2d.NewGeometry_Origin()
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat4  proj;		// Projection matrix
		uniform   mat4  vwmd;		// View*Model matrix
		uniform   vec2  asp;		// aspect ratio, w : h
		uniform   vec3  orgn;		// origin of the label (WORLD XYZ coordinates)
		uniform   vec3  offr;		// offset (CAMERA XY in pixels) & rotation angle
		uniform   vec3  whlen;		// character width & height, and label length
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		attribute vec2  cpose;		// character index & code
		varying float v_code; 		// character code (index of the character in the alphabet texture)
		void main() {
			vec4 origin = proj * vwmd * vec4(orgn, 1.0);
			if (origin.w != 0.0) { origin = origin / origin.w; }
			vec2 ch_off = vec2(offr.x + whlen[0]/2.0, offr.y) + vec2(cpose[0] * whlen[0], 0.0);
			vec2 offset = vec2((ch_off.x + gvxy.x) * 2.0 / asp[0], (ch_off.y + gvxy.y) * 2.0 / asp[1]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
			gl_PointSize = whlen[1];	// character height
			v_code  = cpose[1];
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
		uniform   vec3 	whlen;		// character width & height, and label length
		uniform   vec4  color;		// color of the label
		varying   float v_code;     // character code (index of the character in the alphabet texture)
		void main() {
			vec2 uv = gl_PointCoord;
			if (uv[0] < 0.0 || uv[0] > 1.0) discard;
			if (uv[1] < 0.0 || uv[1] > 1.0) discard;
			float u = uv[0] * (whlen[1]/whlen[0]) - 0.5, v = uv[1];
			if (u < 0.0 || u > 1.0 || v < 0.0 || v > 1.0) discard;
			uv = vec2((u + v_code)/whlen[2], v);	// position in the texture (relative to alphabet_length)
			gl_FragColor = texture2D(texture, uv) * color;
		}`
	offr := []float32{float32(self.offset[0]), float32(self.offset[1]), 0}
	whlen := []float32{self.chwh[0], self.chwh[1], float32(alphabet_texture.GetAlaphabetLength())}
	lrgba := wcommon.GetRGBAFromString(self.color) // label color RGBA
	shader, _ := wcommon.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform("proj", "mat4", "renderer.proj")            // Projection matrix
	shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")            // View*Model matrix
	shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
	shader.SetBindingForUniform("orgn", "vec3", self.xyz[:])                // label origin
	shader.SetBindingForUniform("offr", "vec3", offr)                       // label offset & rotation_angle
	shader.SetBindingForUniform("whlen", "vec3", whlen)                     // ch_width, ch_height, alphabet_length
	shader.SetBindingForUniform("color", "vec4", lrgba[:])                  // label color
	shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
	shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")        // point coordinates
	shader.SetBindingForAttribute("cpose", "vec2", "instance.pose:2:0")     // character pose (:<stride>:<offset>)
	shader.CheckBindings()                                                  // check validity of the shader
	scnobj := NewSceneObject(geometry, alphabet_texture, shader, nil, nil)  // shader for drawing POINTS (for each character)
	scnobj.SetPoses(alphabet_texture.GetAlaphabetPosesForLabel(self.text))
	scnobj.UseDepth = true
	scnobj.UseBlend = true
	self.txtobj = scnobj
	return self
}

func (self *OverlayLabel) build_background_object(wctx *wcommon.WebGLContext) *OverlayLabel {
	tlen := self.chwh[0] * float32(len([]rune(self.text)))
	ltop := [2]float32{self.offset[0] - 4, self.offset[1] + self.chwh[1]/2}
	lbtm := [2]float32{self.offset[0] - 4, self.offset[1] - self.chwh[1]/2}
	rtop := [2]float32{self.offset[0] + tlen + 4, self.offset[1] + self.chwh[1]/2}
	rbtm := [2]float32{self.offset[0] + tlen + 4, self.offset[1] - self.chwh[1]/2}
	bkgtype_split := strings.Split(self.bkgtype, ":")
	if len(bkgtype_split) < 2 {
		fmt.Printf("Failed to build_background_object() : invalid background type '%s'\n", self.bkgtype)
		return self
	}
	bkgtype0 := bkgtype_split[0]
	geometry := webgl2d.NewGeometry()
	material := wcommon.NewMaterial(wctx, bkgtype_split[1])
	switch bkgtype0 {
	case "box": // "box:#ffff00:#000000", "box:<FillColor>:<BorderColor>"
		geometry.SetVertices([][2]float32{lbtm, rbtm, rtop, ltop})
		geometry.SetEdges([][]uint32{{0, 1, 2, 3, 0}})
		geometry.SetFaces([][]uint32{{0, 1, 2, 3}})
		geometry.BuildDataBuffers(true, true, true)
		if len(bkgtype_split) >= 3 { // EDGE color (border color)
			material.SetColorForDrawMode(2, bkgtype_split[2])
		}
	case "under": // "under:#000000", "under:<UnderlineColor>"
		geometry.SetVertices([][2]float32{{0, 0}, lbtm, rbtm})
		switch self.offref {
		case "L_TOP", "L_MID", "L_BTM":
			geometry.SetEdges([][]uint32{{0, 1, 2}})
		case "R_TOP", "R_MID", "R_BTM":
			geometry.SetEdges([][]uint32{{1, 2, 0}})
		case "M_TOP", "CENTER", "M_BTM":
			geometry.SetEdges([][]uint32{{1, 2}})
		default:
			geometry.SetEdges([][]uint32{{1, 2}})
		}
		geometry.BuildDataBuffers(true, true, false)
	default:
		fmt.Printf("Failed to build_background_object() : invalid background type '%s'\n", self.bkgtype)
		return self
	}
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat4  proj;		// Projection matrix
		uniform   mat4  vwmd;		// View*Model matrix
		uniform   vec3  orgn;		// origin of the label (WORLD XYZ coordinates)
		uniform   vec2  asp;		// aspect ratio, w : h
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		void main() {
			vec4 origin = proj * vwmd * vec4(orgn, 1.0);
			if (origin.w != 0.0) { origin = origin / origin.w; }
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
			gl_Position = vec4(origin.xy + offset.xy, origin.z, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// color RGBA
		void main() {
			gl_FragColor = color;
		}`
	shader, _ := wcommon.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform("proj", "mat4", "renderer.proj")      // Projection matrix
	shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")      // View*Model matrix
	shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")     // AspectRatio
	shader.SetBindingForUniform("orgn", "vec3", self.xyz[:])          // label origin
	shader.SetBindingForUniform("color", "vec4", "material.color")    // label color
	shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")  // point coordinates from 2D geometry
	shader.CheckBindings()                                            // check validity of the shader
	scnobj := NewSceneObject(geometry, material, nil, shader, shader) // shader for drawing EDGEs & FACEs
	scnobj.UseDepth = true
	self.bkgobj = scnobj
	return self
}
