package webgl2d

import (
	"fmt"
	"strings"

	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/wcommon/geom2d"
)

var alphabet_texture *wcommon.Material = nil // Alphabet texture

type OverlayLabel struct {
	text    string       // text of the label
	xy      [2]float32   // origin of the label (in WORLD space)
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
	wctx   *wcommon.WebGLContext //
	Labels []*OverlayLabel       //
}

func NewOverlayLabelLayer(wctx *wcommon.WebGLContext, outlined bool) *OverlayLabelLayer {
	self := OverlayLabelLayer{wctx: wctx}
	if alphabet_texture == nil {
		alphabet_texture = wcommon.NewMaterial_AlphabetTexture(wctx, 20, "#ffffff", outlined)
	}
	self.Labels = make([]*OverlayLabel, 0)
	return &self
}

func (self *OverlayLabelLayer) Render(pvm *geom2d.Matrix3) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.wctx)
	for _, label := range self.Labels {
		if label.bkgobj != nil {
			renderer.RenderSceneObject(label.bkgobj, pvm)
		}
		if label.txtobj != nil {
			renderer.RenderSceneObject(label.txtobj, pvm)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Labels
// ----------------------------------------------------------------------------

func (self *OverlayLabelLayer) CreateLabel(label_text string, xy [2]float32, color string) *OverlayLabel {
	chwh := alphabet_texture.GetAlaphabetCharacterWH(1.0)
	return &OverlayLabel{text: label_text, xy: xy, chwh: chwh, color: color}
}

func (self *OverlayLabelLayer) FindLabel(label_text string) *OverlayLabel {
	for _, label := range self.Labels {
		if label.text == label_text {
			return label
		}
	}
	return nil
}

func (self *OverlayLabelLayer) AddLabel(label *OverlayLabel, build bool) *OverlayLabelLayer {
	self.Labels = append(self.Labels, label)
	if build {
		label.BuildSceneObjects(self.wctx)
	}
	return self
}

func (self *OverlayLabelLayer) AddLabelText(label_text string, xy [2]float32, color string, offref string) *OverlayLabelLayer {
	// Simple version to add a Label, instead of doing
	//   label := layer.CreateLabel();  label.SetPose().BuildSceneObjects();  layer.AddLabel(label)
	chwh := alphabet_texture.GetAlaphabetCharacterWH(1.0)
	label := &OverlayLabel{text: label_text, xy: xy, chwh: chwh, color: color}
	label.SetPose(0, offref, [2]float32{0, 0})
	self.Labels = append(self.Labels, label)
	label.build_labeltext_object(self.wctx) // SceneObject for the label text is built (without background)
	return self
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

func (self *OverlayLabel) BuildSceneObjects(wctx *wcommon.WebGLContext) *OverlayLabel {
	self.build_labeltext_object(wctx)  // rebuild SceneObject for text
	self.build_background_object(wctx) // rebuild SceneObject for background
	return self
}

func (self *OverlayLabel) build_labeltext_object(wctx *wcommon.WebGLContext) *OverlayLabel {
	geometry := NewGeometry()
	geometry.SetVertices([][2]float32{{0, 0}}) // geometry with a single vertex
	geometry.BuildDataBuffers(true, false, false)
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3  pvm;		// Projection * View * Model matrix
		uniform   vec2  asp;		// aspect ratio, w : h
		uniform   vec4  orgoff;		// origin & offset of the label (origin: WORLD_XY, offset: CAMERA XY in pixel)
		uniform   vec4  whrlen;		// character width & height, and label length
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		attribute vec2  chic;		// character index & code for each character
		varying float v_code; 		// index of the character in the alphabet texture
		void main() {
			vec3 origin = pvm * vec3(orgoff.xy, 1.0);
			vec2 lb_off = vec2( orgoff.z + whrlen[0]/2.0, orgoff.w );
			vec2 ch_off = vec2( chic[0] * whrlen[0], 0.0 );
			vec2 offset = (lb_off + ch_off + gvxy) * 2.0 / asp[0];
			// vec2 offset = vec2(offset.x * 2.0 / asp[0], offset.y * 2.0 / asp[0]);
			gl_Position = vec4(origin.xy + offset.xy, 0.0, 1.0);
			gl_PointSize = whrlen[1];	// character height
			v_code  = chic[1];
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
		uniform   vec4 	whrlen;		// character width & height, and label length
		uniform   vec4  color;		// color of the label
		varying   float v_code;     // character code (index of the alphabet texture)
		void main() {
			vec2 uv = gl_PointCoord;
			if (uv[0] < 0.0 || uv[0] > 1.0) discard;
			if (uv[1] < 0.0 || uv[1] > 1.0) discard;
			float u = uv[0] * (whrlen[1]/whrlen[0]) - 0.5, v = uv[1];
			if (u < 0.0 || u > 1.0 || v < 0.0 || v > 1.0) discard;
			uv = vec2((u + v_code)/whrlen[3], v);	// position in the texture (relative to label_length)
			gl_FragColor = texture2D(texture, uv) * color;
		}`
	orgoff := []float32{self.xy[0], self.xy[1], float32(self.offset[0]), float32(self.offset[1])}
	whrlen := []float32{self.chwh[0], self.chwh[1], 0, float32(alphabet_texture.GetAlaphabetLength())}
	lbrgba := wcommon.GetRGBAFromString(self.color) // label color RGBA
	// fmt.Println(orgoff, whrlen, lbrgba)
	shader, _ := wcommon.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")              // Proj*View*Model matrix
	shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
	shader.SetBindingForUniform("orgoff", "vec4", orgoff)                   // label origin & offset
	shader.SetBindingForUniform("whrlen", "vec4", whrlen)                   // ch_width, ch_height, rotation, label_length
	shader.SetBindingForUniform("color", "vec4", lbrgba[:])                 // label color
	shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
	shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")        // point coordinates
	shader.SetBindingForAttribute("chic", "vec2", "instance.pose:2:0")      // instance pose (:<stride>:<offset>)
	shader.CheckBindings()                                                  // check validity of the shader
	scnobj := NewSceneObject(geometry, alphabet_texture, shader, nil, nil)  // shader for drawing POINTS (for each character)
	label_rune := []rune(self.text)
	charcodes := wcommon.NewSceneObjectPoses(2, len(label_rune), nil)
	for i := 0; i < len(label_rune); i++ { // save character index & code for each rune
		charcodes.SetPose(i, 0, float32(i), float32(alphabet_texture.GetAlaphabetCharacterIndex(label_rune[i])))
	}
	// fmt.Printf("%v  %v\n", label_rune, cidx_list)
	scnobj.SetInstancePoses(charcodes) // SceneObjectPoses
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
		fmt.Printf("Failed to BuildBkgObject() : invalid background type '%s'\n", self.bkgtype)
		return self
	}
	bkgtype0 := bkgtype_split[0]
	geometry := NewGeometry()
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
		fmt.Printf("Failed to BuildBkgObject() : invalid background type '%s'\n", self.bkgtype)
		return self
	}
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3  pvm;		// Projection * View * Model matrix
		uniform   vec2  org;		// origin of the label (WORLD XY coordinates)
		uniform   vec2  asp;		// aspect ratio, w : h
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		void main() {
			vec3 origin = pvm * vec3(org, 1.0);
			vec2 offset = gvxy * 2.0 / asp[0];
			gl_Position = vec4(origin.xy + offset.xy, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// color RGBA
		void main() {
			gl_FragColor = color;
		}`
	shader, _ := wcommon.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform("pvm", "mat3", "renderer.pvm")        // Proj*View*Model matrix
	shader.SetBindingForUniform("org", "vec2", self.xy[:])            // label origin
	shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")     // AspectRatio
	shader.SetBindingForUniform("color", "vec4", "material.color")    // label color
	shader.SetBindingForAttribute("gvxy", "vec2", "geometry.coords")  // point coordinates
	shader.CheckBindings()                                            // check validity of the shader
	scnobj := NewSceneObject(geometry, material, nil, shader, shader) // shader for drawing EDGEs & FACEs
	self.bkgobj = scnobj
	return self
}
