package webgl3d

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
	texture    js.Value             // TODO: NOT IMPLEMENTED YET
	texture_wh [2]int               // texture size
}

func NewMaterial(wctx *common.WebGLContext, color_or_texture string) *Material {
	mat := Material{wctx: wctx}
	mat.texture = js.Null()
	if color_or_texture[0] == '#' { // simple color
		mat.color, _ = common.ParseHexColor(color_or_texture)
		mat.LoadTexture("") // texture image with a single blue pixel
	} else { // texture
		mat.color = [4]float32{1.0, 1.0, 1.0, 1.0} // white
		mat.LoadTexture(color_or_texture)          // texture image from server path
	}
	return &mat
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

func (self *Material) LoadTexture(path string) {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	if self.texture.IsNull() { // initialize it with (1x1) image of a single blue pixel
		self.texture = context.Call("createTexture")
		context.Call("bindTexture", constants.TEXTURE_2D, self.texture)
		context.Call("texImage2D", constants.TEXTURE_2D, 0, constants.RGBA, 1, 1, 0, constants.RGBA, constants.UNSIGNED_BYTE,
			common.ConvertGoSliceToJsTypedArray([]uint8{0, 0, 255, 255}))
		// gl.texImage2D(gl.TEXTURE_2D, level, internalFormat, width, height, border, srcFormat, srcType, pixels);
		self.texture_wh = [2]int{1, 1}
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
					size := image.Bounds().Size()
					pbuffer := make([]uint8, size.X*size.Y*4)
					for y := 0; y < size.Y; y++ {
						y_idx := (size.Y - 1 - y) * size.X * 4
						for x := 0; x < size.X; x++ {
							rgba := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
							idx := y_idx + x*4
							pbuffer[idx+0] = rgba.R
							pbuffer[idx+1] = rgba.G
							pbuffer[idx+2] = rgba.B
							pbuffer[idx+3] = rgba.A
						}
					}
					context.Call("bindTexture", constants.TEXTURE_2D, self.texture)
					context.Call("texImage2D", constants.TEXTURE_2D, 0, constants.RGBA, size.X, size.Y, 0, constants.RGBA, constants.UNSIGNED_BYTE,
						common.ConvertGoSliceToJsTypedArray(pbuffer))
					is_power_of_two := func(n int) bool {
						return (n & (n - 1)) == 0
					}
					if is_power_of_two(size.X) && is_power_of_two(size.Y) {
						//gl.generateMipmap(gl.TEXTURE_2D);
						context.Call("generateMipmap", constants.TEXTURE_2D)
					} else { // NON-POWER-OF-2 textures
						// WebGL1 can only use FILTERING == NEAREST or LINEAR, and WRAPPING_MODE == CLAMP_TO_EDGE
						// gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
						// gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
						// gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
						context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_WRAP_S, constants.CLAMP_TO_EDGE)
						context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_WRAP_T, constants.CLAMP_TO_EDGE)
						context.Call("texParameteri", constants.TEXTURE_2D, constants.TEXTURE_MIN_FILTER, constants.LINEAR)
					}
					self.texture_wh = [2]int{size.X, size.Y}
				}
			}
		}()
	}
}

func (self *Material) ShowInfo() {
	c := self.color
	fmt.Printf("Material with COLOR [%.2f %.2f %.2f %.2f] and TEXTURE %d x %d\n",
		c[0], c[1], c[2], c[3], self.texture_wh[0], self.texture_wh[1])
}
