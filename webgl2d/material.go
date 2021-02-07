package webgl2d

import "fmt"

type Material struct {
	color   []float32 //
	texture string    // TODO: NOT IMPLEMENTED YET
}

func NewMaterial(color_or_texture string) *Material {
	var mat Material
	if color_or_texture[0] == '#' {
		mat.color, _ = parse_hex_color(color_or_texture)
		mat.texture = ""
	} else {
		mat.color = []float32{1.0, 1.0, 1.0, 1.0} // white
		mat.texture = color_or_texture
	}
	return &mat
}

func parse_hex_color(s string) ([]float32, error) {
	color := []float32{1.0, 1.0, 1.0, 1.0}
	var err error = nil
	if s[0] == '#' {
		c := []uint8{255, 255, 255, 255}
		switch len(s) {
		case 9:
			_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c[0], &c[1], &c[2], &c[3])
		case 7:
			_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c[0], &c[1], &c[2])
		case 5:
			_, err = fmt.Sscanf(s, "#%1x%1x%1x%1x", &c[0], &c[1], &c[2], &c[3])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
			c[3] *= 17
		case 4:
			_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c[0], &c[1], &c[2])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
		default:
			err = fmt.Errorf("invalid length, must be 7 or 4")
		}
		color = []float32{float32(c[0]) / 255.0, float32(c[1]) / 255.0, float32(c[2]) / 255.0, float32(c[3]) / 255.0}
	}
	return color, err
}

func (self *Material) ShowInfo() {
	fmt.Printf("Material with color %v\n", self.color)
}
