package webgl2d

import "github.com/go4orward/gowebgl/common"

type Camera struct {
	ViewMatrix *Matrix3
}

func NewCamera(wctx *common.WebGLContext) *Camera {
	var camera Camera
	return &camera
}

func (self *Camera) SetPose(cx float32, cy float32, width float32, height float32) *Camera {
	return self
}
