package webgl2d

import (
	"github.com/go4orward/gowebgl/geom2d"
)

type Overlay interface {
	Render(pvm *geom2d.Matrix3)
}
