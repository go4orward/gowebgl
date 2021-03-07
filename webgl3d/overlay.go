package webgl3d

import "github.com/go4orward/gowebgl/wcommon/geom3d"

type Overlay interface {
	Render(proj *geom3d.Matrix4, view *geom3d.Matrix4)
}
