package webgl2d

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
)

type Material struct {
	wctx       *common.WebGLContext //
	color      [4]float32           //
	texture    js.Value             // texture
	texture_wh [2]int               // texture size
}

func NewMaterial(wctx *common.WebGLContext, source string) *Material {
	mat := Material{wctx: wctx, color: [4]float32{1, 1, 0, 1}, texture: js.Null(), texture_wh: [2]int{0, 0}}
	if len(source) > 0 {
		if source[0] == '#' { // COLOR RGB value
			mat.SetColor(source)
		} else if source[0] == '/' { // TEXTURE image path
			mat.LoadTexture(source)
		}
	}
	return &mat
}

func (self *Material) SetColor(color string) *Material {
	self.color, _ = common.ParseHexColor(color) // 'color' like "#ff0000ff"
	return self
}

func (self *Material) GetColor() [4]float32 {
	return self.color
}

func (self *Material) GetTexture() js.Value {
	return self.texture
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *Material) LoadTextureOfSingleBluePixel() *Material {
	// Load texture from 1x1 image of a single blue pixel
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	if self.texture.IsNull() {
		self.texture = context.Call("createTexture")
	}
	context.Call("bindTexture", constants.TEXTURE_2D, self.texture)
	context.Call("texImage2D", constants.TEXTURE_2D, 0, constants.RGBA, 1, 1, 0, constants.RGBA, constants.UNSIGNED_BYTE,
		common.ConvertGoSliceToJsTypedArray([]uint8{0, 0, 255, 255}))
	// gl.texImage2D(gl.TEXTURE_2D, level, internalFormat, width, height, border, srcFormat, srcType, pixels);
	self.texture_wh = [2]int{1, 1}
	return self
}

func (self *Material) LoadTexture(path string) *Material {
	// Load texture image from server path, for example "/assets/world.jpg"
	if self.texture.IsNull() {
		self.LoadTextureOfSingleBluePixel() // initialize it with a single blue pixel
	}
	if path != "" {
		go func() {
			resp, err := http.Get(path)
			if err != nil {
				fmt.Printf("Failed to GET %s : %v\n", path, err)
				return
			}
			defer resp.Body.Close()
			response_body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to GET %s : %v\n", path, err)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 { // response with error message
				fmt.Printf("Failed to GET %s : (%d) %s\n", path, resp.StatusCode, string(response_body))
			} else { // successful response with texture image
				var image image.Image
				var err error
				// log.Printf("Texture image downloaded from server\n")
				switch filepath.Ext(path) {
				case ".png", ".PNG":
					image, err = png.Decode(bytes.NewBuffer(response_body))
				case ".jpg", ".JPG":
					image, err = jpeg.Decode(bytes.NewBuffer(response_body))
				default:
					fmt.Printf("Invalid image format for %s\n", path)
					return
				}
				if err != nil {
					fmt.Printf("Failed to decode %s : %v\n", path, err)
				} else {
					// log.Printf("Texture image decoded (by png/jpg)\n")
					size := image.Bounds().Size()
					pbuffer := make([]uint8, size.X*size.Y*4)
					for y := 0; y < size.Y; y++ {
						y_idx := (size.Y - 1 - y) * size.X * 4
						for x := 0; x < size.X; x++ {
							rgba := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
							idx := y_idx + x*4
							set_pbuffer_with_rgba(pbuffer, idx, rgba.R, rgba.G, rgba.B, rgba.A)
						}
					}
					// log.Printf("Texture pixel buffer converted to RGBA\n")
					self.LoadTextureFromBufferRGBA(pbuffer, size.X, size.Y)
					// log.Printf("Texture ready for WebGL\n")
				}
			}
		}()
	}
	return self
}

func (self *Material) LoadTextureFromBufferRGBA(buffer []uint8, width int, height int) *Material {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	js_buffer := common.ConvertGoSliceToJsTypedArray(buffer)
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
		context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_MIN_FILTER, constants.LINEAR)
	}
	self.texture_wh = [2]int{width, height}
	return self
}

func (self *Material) ShowInfo() {
	c := self.color
	fmt.Printf("Material with COLOR [%.2f %.2f %.2f %.2f] and TEXTURE %d x %d\n",
		c[0], c[1], c[2], c[3], self.texture_wh[0], self.texture_wh[1])
}

func set_pbuffer_with_rgba(pbuffer []uint8, idx int, R uint8, G uint8, B uint8, A uint8) {
	pbuffer[idx+0] = R
	pbuffer[idx+1] = G
	pbuffer[idx+2] = B
	pbuffer[idx+3] = A
}
