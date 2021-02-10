package webgl2d

import (
	"fmt"

	"github.com/go4orward/gowebgl/common"
)

type Material struct {
	color   [4]float32 //
	texture string     // TODO: NOT IMPLEMENTED YET
}

func NewMaterial(color_or_texture string) *Material {
	var mat Material
	if color_or_texture[0] == '#' {
		mat.color, _ = common.ParseHexColor(color_or_texture)
		mat.texture = ""
	} else {
		mat.color = [4]float32{1.0, 1.0, 1.0, 1.0} // white
		mat.texture = color_or_texture
	}
	return &mat
}

func (self *Material) ShowInfo() {
	fmt.Printf("Material with color %v\n", self.color)
}
