package wcommon

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"syscall/js"
)

type Material struct {
	wctx            *WebGLContext //
	color           [4][4]float32 // color ([0]:common, [1]:vert, [2]:edge, [3]:face)
	texture         js.Value      // texture
	texture_wh      [2]int        // texture size
	alphabet_cwh    [2]float32    // character width & height of ALPHABET_STRING
	texture_loading bool          // true, only if texture is being loaded
}

func NewMaterial(wctx *WebGLContext, source string) *Material {
	mat := Material{wctx: wctx, texture: js.Null(), texture_wh: [2]int{0, 0}}
	mat.SetDrawModeColor(0, [4]float32{0, 1, 1, 1})
	if len(source) > 0 {
		if source[0] == '#' { // COLOR RGB value
			rgba := GetRGBAFromString(source)
			mat.SetDrawModeColor(0, rgba)
		} else { // TEXTURE image path
			mat.LoadTexture(source)
		}
	}
	return &mat
}

func (self *Material) ShowInfo() {
	colors := ""
	for i := 0; i < len(self.color); i++ {
		c := [4]uint8{uint8(self.color[i][0] * 255), uint8(self.color[i][1] * 255), uint8(self.color[i][2] * 255), uint8(self.color[i][3] * 255)}
		colors += fmt.Sprintf("#%02x%02x%02x%02x ", c[0], c[1], c[2], c[3])
	}
	fmt.Printf("Material with TEXTURE %dx%d and COLOR %s\n", self.texture_wh[0], self.texture_wh[1], colors)
}

// ----------------------------------------------------------------------------
// COLOR
// ----------------------------------------------------------------------------

func (self *Material) SetColorForDrawMode(draw_mode int, color string) *Material {
	// 'draw_mode' :  0:common, 1:vertex, 2:edges, 3:faces
	if color != "" {
		return self.SetDrawModeColor(draw_mode, GetRGBAFromString(color))
	} else {
		return self
	}
}

func (self *Material) SetDrawModeColor(draw_mode int, color [4]float32) *Material {
	switch draw_mode {
	case 1:
		self.color[1] = color // vertex color
	case 2:
		self.color[2] = color // edge color
	case 3:
		self.color[3] = color // face color
	default:
		self.color[0] = color // otherwise
		self.color[1] = color
		self.color[2] = color
		self.color[3] = color
	}
	return self
}

func (self *Material) GetDrawModeColor(draw_mode int) [4]float32 {
	return self.color[draw_mode]
}

func GetRGBAFromString(s string) [4]float32 {
	c := [4]uint8{0, 0, 0, 255}
	if len(s) == 0 {
	} else if s[0] == '#' {
		switch len(s) {
		case 9:
			fmt.Sscanf(s, "#%02x%02x%02x%02x", &c[0], &c[1], &c[2], &c[3])
		case 7:
			fmt.Sscanf(s, "#%02x%02x%02x", &c[0], &c[1], &c[2])
		case 5:
			fmt.Sscanf(s, "#%1x%1x%1x%1x", &c[0], &c[1], &c[2], &c[3])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
			c[3] *= 17
		case 4:
			fmt.Sscanf(s, "#%1x%1x%1x", &c[0], &c[1], &c[2])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
		default:
		}
	} else {
	}
	return [4]float32{float32(c[0]) / 255, float32(c[1]) / 255, float32(c[2]) / 255, float32(c[3]) / 255}
}

// ----------------------------------------------------------------------------
// TEXTURE
// ----------------------------------------------------------------------------

func (self *Material) GetTexture() js.Value {
	return self.texture
}

func (self *Material) GetTextureWH() [2]int {
	return self.texture_wh
}

func (self *Material) IsTextureReady() bool {
	return (!self.texture.IsNull() && self.texture_wh[0] > 0 && self.texture_wh[1] > 0)
}

func (self *Material) IsTextureLoading() bool {
	return self.texture_loading
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *Material) LoadTexture(path string) *Material {
	// Load texture image from server path, for example "/assets/world.jpg"
	context, c := self.wctx.GetContext(), self.wctx.GetConstants()
	if self.texture.IsNull() { // initialize it with a single CYAN pixel
		self.texture = context.Call("createTexture")
		context.Call("bindTexture", c.TEXTURE_2D, self.texture)
		context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, 1, 1, 0, c.RGBA, c.UNSIGNED_BYTE, ConvertGoSliceToJsTypedArray([]uint8{0, 255, 255, 255}))
		self.texture_wh = [2]int{1, 1}
	}
	self.texture_loading = true
	if path != "" {
		go func() {
			defer func() { self.texture_loading = false }()
			// log.Printf("Texture started GET %s\n", path)
			resp, err := http.Get(path)
			if err != nil {
				log.Printf("Failed to GET %s : %v\n", path, err)
				return
			}
			defer resp.Body.Close()
			response_body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to GET %s : %v\n", path, err)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 { // response with error message
				log.Printf("Failed to GET %s : (%d) %s\n", path, resp.StatusCode, string(response_body))
			} else { // successful response with texture image
				// log.Printf("Texture image downloaded from server\n")
				var img image.Image
				var err error
				switch filepath.Ext(path) {
				case ".png", ".PNG":
					img, err = png.Decode(bytes.NewBuffer(response_body))
				case ".jpg", ".JPG":
					img, err = jpeg.Decode(bytes.NewBuffer(response_body))
				default:
					fmt.Printf("Invalid image format for %s\n", path)
					return
				}
				if err != nil {
					log.Printf("Failed to decode %s : %v\n", path, err)
				} else {
					size := img.Bounds().Size()
					// log.Printf("Texture image (%dx%d) decoded as %T\n", size.X, size.Y, img)
					var pixbuf []uint8
					switch img.(type) {
					case *image.RGBA: // traditional 32-bit alpha-premultiplied R/G/B/A color
						pixbuf = img.(*image.RGBA).Pix
					case *image.NRGBA: // non-alpha-premultiplied 32-bit R/G/B/A color
						pixbuf = img.(*image.NRGBA).Pix
					default: // we need conversion, otherwise
						pixbuf = make([]uint8, size.X*size.Y*4)
						for y := 0; y < size.Y; y++ {
							y_idx := y * size.X * 4
							for x := 0; x < size.X; x++ {
								rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
								idx := y_idx + x*4
								set_pixbuf_with_rgba(pixbuf, idx, rgba.R, rgba.G, rgba.B, rgba.A)
							}
						}
					}
					js_buffer := ConvertGoSliceToJsTypedArray(pixbuf)
					context.Call("bindTexture", c.TEXTURE_2D, self.texture)
					context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, size.X, size.Y, 0, c.RGBA, c.UNSIGNED_BYTE, js_buffer)
					if size.X&(size.X-1) == 0 && size.Y&(size.Y-1) == 0 { // POWER-OF-2 width & height
						context.Call("generateMipmap", c.TEXTURE_2D)
					} else { // NON-POWER-OF-2 textures : CLAMP_TO_EDGE & NEAREST/LINEAR only
						context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_S, c.CLAMP_TO_EDGE)
						context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_T, c.CLAMP_TO_EDGE)
						context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, c.LINEAR)
					}
					self.texture_wh = [2]int{size.X, size.Y}
					// log.Printf("Texture ready for WebGL\n")
				}
			}
		}()
	}
	return self
}

func set_pixbuf_with_rgba(pbuffer []uint8, idx int, R uint8, G uint8, B uint8, A uint8) {
	pbuffer[idx+0] = R
	pbuffer[idx+1] = G
	pbuffer[idx+2] = B
	pbuffer[idx+3] = A
}

// ----------------------------------------------------------------------------
// Glow Texture
// ----------------------------------------------------------------------------

func NewMaterial_GlowTexture(wctx *WebGLContext, color string) *Material {
	self := Material{wctx: wctx, texture: js.Null(), texture_wh: [2]int{0, 0}}
	// decide RGB color
	rgba := [4]float32{0, 1, 1, 1}
	if len(color) > 0 && color[0] == '#' {
		rgba = GetRGBAFromString(color)
	}
	self.SetDrawModeColor(0, rgba)
	// build pixel buffer for glow effect
	const width, height = 34, 2 // it has to be non-power-of-two texture with gl.NEAREST
	pixbuf := make([]uint8, (width*height)*4)
	for u := 1; u < width-1; u++ { // first (i==0) and last (i==width-1) pixel is ZERO
		ratio := (float32(u-1) / float32(width-2))
		if true { // diminishing glow for the first row (v == 0)  [ 1.0 ~ 0.5 ~ 0.0 ]
			intensity := 1.0 - ratio
			ii := intensity * intensity
			set_pixbuf_with_rgba(pixbuf, (u)*4, uint8(ii*rgba[0]*255), uint8(ii*rgba[1]*255), uint8(ii*rgba[2]*255), uint8(ii*255))
		}
		if true { // glow on both side for the second row (v == 1)  [ 0.0 ~ 1.0 ~ 0.0 ]
			intensity := 1.0 - float32(math.Abs(float64(ratio*2-1)))
			ii := intensity * intensity
			set_pixbuf_with_rgba(pixbuf, (width+u)*4, uint8(ii*rgba[0]*255), uint8(ii*rgba[1]*255), uint8(ii*rgba[2]*255), uint8(ii*255))
		}
	}
	context, c := self.wctx.GetContext(), self.wctx.GetConstants()
	js_buffer := ConvertGoSliceToJsTypedArray(pixbuf)
	self.texture = context.Call("createTexture")
	context.Call("bindTexture", c.TEXTURE_2D, self.texture)
	context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, width, height, 0, c.RGBA, c.UNSIGNED_BYTE, js_buffer)
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_S, c.CLAMP_TO_EDGE)
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_T, c.CLAMP_TO_EDGE)
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, c.NEAREST)
	self.texture_wh = [2]int{width, height} // CLAMP_TO_EDGE & NEAREST(not LINEAR) for NON-POWER-OF-2 textures
	return &self
}

// ----------------------------------------------------------------------------
// Alphabet Texture  (for labels)
// ----------------------------------------------------------------------------

const _ALPHABET_STRING = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~?°"

func NewMaterial_AlphabetTexture(wctx *WebGLContext, fontsize int, color string, outlined bool) *Material {
	self := Material{wctx: wctx, texture: js.Null(), texture_wh: [2]int{0, 0}}
	self.SetColorForDrawMode(0, color)
	// 'fontsize' : 12=>(7.2x12.6), 16=>(9.6x16.8), 20=>(12x21), 24=>(14x25), 30=>(18x31), 40=>(24x42)
	fontstyle := "Courier New" // "Courier" or "Monospace" or "Courier New"
	_ALPHABET_FONT_STYLE := fmt.Sprintf("%dpx %s", fontsize, fontstyle)
	txtctx := js.Global().Get("document").Call("createElement", "canvas").Call("getContext", "2d")
	txtctx.Set("font", _ALPHABET_FONT_STYLE) // need to be set, before measuring text size
	cwidth := float32(txtctx.Call("measureText", "M").Get("width").Float())
	cheight := float32(fontsize) * 1.05 // we need some more margin below the text
	twidth := int(math.Floor(txtctx.Call("measureText", _ALPHABET_STRING).Get("width").Float()))
	theight := int(cheight) // instead of int(cwidth*2)
	// fmt.Printf("Character: %v %v  Texture: %v %v\n", cwidth, cheight, twidth, theight)
	txtctx.Get("canvas").Set("width", twidth)
	txtctx.Get("canvas").Set("height", theight)
	txtctx.Call("clearRect", 0, 0, twidth, theight)
	txtctx.Set("font", _ALPHABET_FONT_STYLE) // need to be set again!
	txtctx.Set("textAlign", "start")         // start (default), end, left, right, center
	txtctx.Set("textBaseline", "top")        // top, hanging, middle, alphabetic (default), ideographic, bottom
	if outlined {
		txtctx.Set("strokeStyle", "#000000")              // BLACK outline
		txtctx.Set("lineWidth", 2.5)                      // text stroke width
		txtctx.Call("strokeText", _ALPHABET_STRING, 0, 0) // draw the alphabet string for outline
	}
	txtctx.Set("fillStyle", color)                  // interior (Note that WHITE can be multiplied with other colors later)
	txtctx.Call("fillText", _ALPHABET_STRING, 0, 0) // draw the alphabet string
	context, c := self.wctx.GetContext(), self.wctx.GetConstants()
	self.texture = context.Call("createTexture")
	context.Call("bindTexture", c.TEXTURE_2D, self.texture)
	context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, c.RGBA, c.UNSIGNED_BYTE, txtctx.Get("canvas"))
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_S, c.CLAMP_TO_EDGE)
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_WRAP_T, c.CLAMP_TO_EDGE)
	context.Call("texParameteri", c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, c.LINEAR)
	self.texture_wh = [2]int{twidth, theight}
	self.alphabet_cwh = [2]float32{cwidth, cheight}
	return &self
}

func (self *Material) GetAlaphabetLength() int {
	return len([]rune(_ALPHABET_STRING)) // total rune length of _ALPHABET_STRING
}

func (self *Material) GetAlaphabetCharacterWH(scale float32) [2]float32 {
	return [2]float32{self.alphabet_cwh[0] * scale, self.alphabet_cwh[1] * scale}
}

func (self *Material) GetAlaphabetCharacterIndex(c rune) int {
	if c >= 32 && c < 127 { // 32: <SPACE>, 127: <DEL>
		return int(c) - 32
	} else {
		alphabet_rune := []rune(_ALPHABET_STRING)
		for i := 127 - 32; i < len(alphabet_rune); i++ {
			if c == alphabet_rune[i] {
				return i
			}
		}
		return (127 - 32) // '?'
	}
}

func (self *Material) GetAlaphabetPosesForLabel(label_text string) *SceneObjectPoses {
	// Get SceneObjectPoses of the 'label_text' (as a list of index & code pairs),
	//   which can be used for rendering the 'label_text' in OverlayLabelLayer
	label_runes := []rune(label_text)
	characters := NewSceneObjectPoses(2, len(label_runes), nil)
	for i := 0; i < len(label_runes); i++ { // save character index & code for each rune
		characters.SetPose(i, 0, float32(i), float32(self.GetAlaphabetCharacterIndex(label_runes[i])))
	}
	return characters
}
