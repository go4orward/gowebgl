package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type Overlay interface {
	Render(wctx *common.WebGLContext, pvm *geom2d.Matrix3)
}
