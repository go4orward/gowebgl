package webgl3d

import (
	"fmt"

	"github.com/go4orward/gowebgl/common/geom3d"
)

type Camera struct {
	center     [3]float32      // camera position in world space
	angle      float32         // camera rotation angle in degree
	zoom       float32         // camera zoom level
	viewmatrix *geom3d.Matrix4 // view matrix   Mcw
}

func NewCamera() *Camera {
	var camera Camera
	camera.viewmatrix = geom3d.NewMatrix4() // identity
	camera.center = [3]float32{0, 0, 0}
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

func (self *Camera) SetPose(center [3]float32, angle_in_degree float32, zoom float32) *Camera {
	if zoom <= 0.0 || zoom >= 1000.0 {
		fmt.Printf("Camera.SetPose() failed : invalid zoom = %.1f\n", zoom)
		return self
	}
	self.center = center
	self.angle = angle_in_degree
	self.zoom = zoom
	return self.update_view_matrix()
}

func (self *Camera) update_view_matrix() *Camera {
	// radian := float64(self.angle) * (math.Pi / 180.0)
	// cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	// scaling := geom3d.NewMatrix4().Set(
	// 	self.zoom, 0.0, 0.0,
	// 	0.0, self.zoom, 0.0,
	// 	0.0, 0.0, 1.0)
	// rotation := geom3d.NewMatrix4().Set(
	// 	cos, +sin, 0.0,
	// 	-sin, cos, 0.0,
	// 	0.0, 0.0, 1.0)
	// translation := geom3d.NewMatrix4().Set(
	// 	1.0, 0.0, -self.center[0],
	// 	0.0, 1.0, -self.center[1],
	// 	0.0, 0.0, 1.0)
	// self.viewmatrix.SetMultiplyMatrices(scaling, rotation, translation)
	return self
}

func (self *Camera) ShowInfo() {
	fmt.Printf("Camera at (%v,%v) with angle=%.1f zoom=%.2f  %v\n", self.center[0], self.center[1], self.angle, self.zoom, self.viewmatrix)
}
