package common

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
	texture_loading bool          // true, only if texture is being loaded
}

func NewMaterial(wctx *WebGLContext, source string) *Material {
	mat := Material{wctx: wctx, texture: js.Null(), texture_wh: [2]int{0, 0}}
	mat.SetDrawModeColor(0, [4]float32{0, 1, 1, 1})
	if len(source) > 0 {
		if source[0] == '#' { // COLOR RGB value
			rgba := parse_hex_color(source)
			mat.SetDrawModeColor(0, rgba)
			mat.LoadTextureOfSinglePixel(rgba)
		} else { // TEXTURE image path
			mat.LoadTexture(source)
		}
	}
	return &mat
}

func NewMaterialForGlowEffect(wctx *WebGLContext, color string) *Material {
	mat := Material{wctx: wctx, texture: js.Null(), texture_wh: [2]int{0, 0}}
	if len(color) > 0 {
		if color[0] == '#' { // COLOR RGB value
			rgba := parse_hex_color(color)
			mat.SetDrawModeColor(0, rgba)
			// Load texture for glow effect
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
			mat.LoadTextureFromBufferRGBA(pixbuf, width, height, false) // non-power-of-two texture with gl.NEAREST
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
	return self.SetDrawModeColor(draw_mode, parse_hex_color(color))
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

func parse_hex_color(s string) [4]float32 {
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

func (self *Material) IsTextureReady() bool {
	return (self.texture_wh[0] > 0 && self.texture_wh[1] > 0)
}

func (self *Material) IsTextureLoading() bool {
	return self.texture_loading
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *Material) LoadTextureOfSinglePixel(rgba [4]float32) *Material {
	// Load texture from 1x1 image of a single pixel
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	if self.texture.IsNull() {
		self.texture = context.Call("createTexture")
	}
	color := []uint8{uint8(rgba[0] * 255), uint8(rgba[1] * 255), uint8(rgba[2] * 255), uint8(rgba[3] * 255)}
	context.Call("bindTexture", constants.TEXTURE_2D, self.texture)
	context.Call("texImage2D", constants.TEXTURE_2D, 0, constants.RGBA, 1, 1, 0, constants.RGBA, constants.UNSIGNED_BYTE,
		ConvertGoSliceToJsTypedArray(color))
	// gl.texImage2D(gl.TEXTURE_2D, level, internalFormat, width, height, border, srcFormat, srcType, pixels);
	self.texture_wh = [2]int{1, 1}
	return self
}

func (self *Material) LoadTexture(path string) *Material {
	// Load texture image from server path, for example "/assets/world.jpg"
	if self.texture.IsNull() { // initialize it with a single CYAN pixel
		self.LoadTextureOfSinglePixel(self.color[0])
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
					default:
						pixbuf = make([]uint8, size.X*size.Y*4)
						for y := 0; y < size.Y; y++ {
							y_idx := y * size.X * 4
							for x := 0; x < size.X; x++ {
								rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
								idx := y_idx + x*4
								set_pixbuf_with_rgba(pixbuf, idx, rgba.R, rgba.G, rgba.B, rgba.A)
							}
						}
						// log.Printf("Texture pixel buffer converted to RGBA\n")
					}
					self.LoadTextureFromBufferRGBA(pixbuf, size.X, size.Y, true)
					// log.Printf("Texture ready for WebGL\n")
				}
			}
		}()
	}
	return self
}

func (self *Material) LoadTextureFromBufferRGBA(buffer []uint8, width int, height int, linear bool) *Material {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	js_buffer := ConvertGoSliceToJsTypedArray(buffer)
	// log.Printf("Texture pixel buffer ready for Javascript\n")
	if self.texture.IsNull() {
		self.texture = context.Call("createTexture")
	}
	context.Call("bindTexture", constants.TEXTURE_2D, self.texture)
	context.Call("texImage2D", constants.TEXTURE_2D, 0, constants.RGBA, width, height, 0, constants.RGBA, constants.UNSIGNED_BYTE, js_buffer)
	is_power_of_two := func(n int) bool {
		return (n & (n - 1)) == 0
	}
	if is_power_of_two(width) && is_power_of_two(height) {
		//gl.generateMipmap(gl.TEXTURE_2D);
		context.Call("generateMipmap", constants.TEXTURE_2D)
	} else { // NON-POWER-OF-2 textures
		// WebGL1 can only use FILTERING == NEAREST or LINEAR, and WRAPPING_MODE == CLAMP_TO_EDGE
		context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_WRAP_S, constants.CLAMP_TO_EDGE)
		context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_WRAP_T, constants.CLAMP_TO_EDGE)
	}
	if linear {
		context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_MIN_FILTER, constants.LINEAR)
	} else {
		context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_MIN_FILTER, constants.NEAREST)
	}
	self.texture_wh = [2]int{width, height}
	return self
}

func set_pixbuf_with_rgba(pbuffer []uint8, idx int, R uint8, G uint8, B uint8, A uint8) {
	pbuffer[idx+0] = R
	pbuffer[idx+1] = G
	pbuffer[idx+2] = B
	pbuffer[idx+3] = A
}
