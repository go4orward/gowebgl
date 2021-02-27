package common

import "fmt"

func ParseHexColor(s string) ([4]uint8, error) {
	c := [4]uint8{255, 255, 255, 255}
	var err error = nil
	if s[0] == '#' {
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
	}
	return c, err
}
