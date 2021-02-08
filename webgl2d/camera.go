package webgl2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type Camera struct {
	center     []float32       // camera position in world space
	angle      float32         // camera rotation angle in degree
	zoom       float32         // camera zoom level
	viewmatrix *geom2d.Matrix3 // view matrix   Mcw
}

func NewCamera(wctx *common.WebGLContext) *Camera {
	var camera Camera
	camera.viewmatrix = geom2d.NewMatrix3() // identity
	camera.center = []float32{0, 0}
	camera.angle = 0.0
	camera.zoom = 1.0
	return &camera
}

func (self *Camera) SetAngle(angle_in_degree float32) *Camera {
	self.angle = angle_in_degree
	return self.update_view_matrix()
}

func (self *Camera) SetCenter(cx float32, cy float32) *Camera {
	self.center[0] = cx
	self.center[1] = cy
	return self.update_view_matrix()
}

func (self *Camera) SetZoom(zoom float32) *Camera {
	self.zoom = zoom
	return self.update_view_matrix()
}

func (self *Camera) SetPose(center []float32, zoom float32, angle_in_degree float32) *Camera {
	geom2d.Assign(self.center, center)
	self.angle = angle_in_degree
	self.zoom = zoom
	return self.update_view_matrix()
}

func (self *Camera) update_view_matrix() *Camera {
	radian := float64(self.angle) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	scale := geom2d.NewMatrix3().Set(
		self.zoom, 0.0, 0.0,
		0.0, self.zoom, 0.0,
		0.0, 0.0, 1.0)
	rotation := geom2d.NewMatrix3().Set(
		cos, -sin, 0.0,
		sin, cos, 0.0,
		0.0, 0.0, 1.0)
	translation := geom2d.NewMatrix3().Set(
		0.0, 0.0, -self.center[0],
		0.0, 0.0, -self.center[1],
		0.0, 0.0, 1.0)
	self.viewmatrix.SetMultiplyMatrices(scale, rotation, translation)
	return self
}

func (self *Camera) ShowInfo() {
	fmt.Printf("Camera at (%v,%v) with angle=%.1f zoom=%.2f\n", self.center[0], self.center[1], self.angle, self.zoom)
	fmt.Println(self.viewmatrix)
}
