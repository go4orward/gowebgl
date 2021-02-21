package webglobe

import (
	"github.com/go4orward/gowebgl/webgl3d"
)

type WorldCamera struct {
	cam3d *webgl3d.Camera
}

func NewWorldCamera(wh [2]int, fov float32, zoom float32, lon float32, lat float32) *WorldCamera {
	cam3d := webgl3d.NewPerspectiveCamera(wh, fov, zoom)
	self := WorldCamera{cam3d: cam3d}
	return &self
}

func (self *WorldCamera) ShowInfo() {
	self.cam3d.ShowInfo()
}
