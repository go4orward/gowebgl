package webglglobe

import (
	"github.com/go4orward/gowebgl/wcommon"
	"github.com/go4orward/gowebgl/webgl2d"
)

type GlobeMarkLayer struct {
	ScnObjs []*webgl2d.SceneObject // 2D SceneObjects to be rendered (in CAMERA space)
}

func NewGlobeMarkLayer() *GlobeMarkLayer {
	mlayer := GlobeMarkLayer{}
	mlayer.ScnObjs = make([]*webgl2d.SceneObject, 0)
	return &mlayer
}

// ----------------------------------------------------------------------------
// Mark 						(single instance with its own geometry)
// ----------------------------------------------------------------------------

func (self *GlobeMarkLayer) AddMark(geometry *webgl2d.Geometry, color string) *GlobeMarkLayer {
	return self
}

// ----------------------------------------------------------------------------
// Mark with Instance Poses 	(multiple instances sharing the same geometry)
// ----------------------------------------------------------------------------

func (self *GlobeMarkLayer) AddMarkWithPoses(geometry *webgl2d.Geometry, poses *wcommon.SceneObjectPoses) *GlobeMarkLayer {
	return self
}
