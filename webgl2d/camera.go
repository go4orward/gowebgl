package webgl2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gowebgl/common/geom2d"
)

type Camera struct {
	// camera internal parameters
	wh         [2]int         // canvas aspect_ratio (width & height)
	fov        float32        // field of view - clip space with (default 2.6)
	zoom       float32        // camera zoom level
	projmatrix geom2d.Matrix3 // projection matrix   Msc (camera space => clip   space)
	// camera pose
	center     [2]float32     // camera position in world space
	viewmatrix geom2d.Matrix3 // view matrix   Mcw (world  space => camera space)
	// final Projection * View matrix
	pjvwmatrix geom2d.Matrix3 //
	// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
	tbbox [2][2]float32 // translation bounding box
}

func NewCamera(wh_aspect_ratio [2]int, fov_in_clipwidth float32, zoom float32) *Camera {
	// FOV 2 covers width 2 in full screen
	camera := Camera{wh: wh_aspect_ratio, fov: fov_in_clipwidth, zoom: zoom}
	camera.update_proj_matrix()
	camera.SetPose(0, 0, 0.0)
	camera.tbbox = geom2d.BBoxInit()
	return &camera
}

func (self *Camera) ShowInfo() {
	p := self.projmatrix.GetElements() // Note that Matrix3 is column-major (just like WebGL)
	v := self.viewmatrix.GetElements()
	fmt.Printf("Camera  centered at [%5.2f %5.2f]\n", self.center[0], self.center[1])
	fmt.Printf("  Parameters : AspectRatio=[%d:%d]  fov=%.1f  zoom=%.2f\n", self.wh[0], self.wh[1], self.fov, self.zoom)
	fmt.Printf("  [ %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %7.2f ]\n", p[0], p[3], p[6], v[0], v[3], v[6])
	fmt.Printf("  [ %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %7.2f ]\n", p[1], p[4], p[7], v[1], v[4], v[7])
	fmt.Printf("  [ %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %7.2f ]\n", p[2], p[5], p[8], v[2], v[5], v[8])
}

// ----------------------------------------------------------------------------
// Camera Internal Parameters
// ----------------------------------------------------------------------------

func (self *Camera) SetAspectRatio(width int, height int) *Camera {
	// This function can be called to handle 'window.resize' event
	self.wh = [2]int{width, height}
	return self.update_proj_matrix()
}

func (self *Camera) SetZoom(zoom float32) *Camera {
	// This function can be called to handle 'wheel' event [ 0.01 ~ 1.0(default) ~ 100.0 ]
	zoom = float32(math.Max(0.001, math.Min(float64(zoom), 1000.0)))
	self.zoom = zoom
	return self.update_proj_matrix()
}

func (self *Camera) update_proj_matrix() *Camera {
	// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
	aspect_ratio := float32(self.wh[0]) / float32(self.wh[1])
	clip_width, clip_height := float32(2.0), 2.0/aspect_ratio // CLIP space width (2.0) & height
	// if aspect_ratio > 1.0 {
	// 	clip_width, clip_height = 2.0*aspect_ratio, float32(2.0) // CLIP space width & height (2.0)
	// }
	x := 2 * self.zoom / clip_width
	y := 2 * self.zoom / clip_height
	ff := 2.0 / self.fov // fov 2.0 will cover width 2 at any distance in full screen
	self.projmatrix.Set(
		ff*x, 0.0, 0.0,
		0.0, ff*y, 0.0,
		0.0, 0.0, 1.0)
	self.pjvwmatrix.SetMultiplyMatrices(&self.projmatrix, &self.viewmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Camera Pose
// ----------------------------------------------------------------------------

func (self *Camera) SetPose(cx float32, cy float32, angle_in_degree float32) *Camera {
	self.center = [2]float32{cx, cy}
	radian := float64(angle_in_degree) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	rotation := geom2d.NewMatrix3().Set(
		cos, +sin, 0.0,
		-sin, cos, 0.0,
		0.0, 0.0, 1.0)
	translation := geom2d.NewMatrix3().Set(
		1.0, 0.0, -self.center[0],
		0.0, 1.0, -self.center[1],
		0.0, 0.0, 1.0)
	self.viewmatrix.SetMultiplyMatrices(rotation, translation)
	self.pjvwmatrix.SetMultiplyMatrices(&self.projmatrix, &self.viewmatrix)
	return self
}

func (self *Camera) Rotate(angle_in_degree float32) *Camera {
	radian := float64(angle_in_degree) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	rotation := geom2d.NewMatrix3().Set(
		cos, +sin, 0.0,
		-sin, cos, 0.0,
		0.0, 0.0, 1.0)
	self.viewmatrix = *rotation.MultiplyRight(&self.viewmatrix)
	self.pjvwmatrix.SetMultiplyMatrices(&self.projmatrix, &self.viewmatrix)
	return self
}

func (self *Camera) Translate(tx float32, ty float32) *Camera {
	if geom2d.BBoxIsSet(self.tbbox) { // translation bounding box
		if (self.center[0] + tx) < self.tbbox[0][0] {
			tx = self.tbbox[0][0] - self.center[0]
		}
		if (self.center[0] + tx) > self.tbbox[1][0] {
			tx = self.tbbox[1][0] - self.center[0]
		}
		if (self.center[1] + ty) < self.tbbox[0][1] {
			ty = self.tbbox[0][1] - self.center[1]
		}
		if (self.center[1] + ty) > self.tbbox[1][1] {
			ty = self.tbbox[1][1] - self.center[1]
		}
	}
	translation := geom2d.NewMatrix3().Set(
		1.0, 0.0, -tx,
		0.0, 1.0, -ty,
		0.0, 0.0, 1.0)
	self.viewmatrix = *translation.MultiplyRight(&self.viewmatrix)
	self.pjvwmatrix.SetMultiplyMatrices(&self.projmatrix, &self.viewmatrix)
	self.center = [2]float32{self.center[0] + tx, self.center[1] + ty}
	return self
}

func (self *Camera) SetTranslationBoundingBox(bbox [2][2]float32) *Camera {
	self.tbbox = bbox
	return self
}

// ----------------------------------------------------------------------------
// Unprojection
// ----------------------------------------------------------------------------

func (self *Camera) UnprojectCanvasToWorld(canvasxy [2]int) [2]float32 {
	hw, hh := (float32(self.wh[0]) / 2), (float32(self.wh[1]) / 2)
	clipxy := [2]float32{(float32(canvasxy[0]) - hw) / hw, -(float32(canvasxy[1]) - hh) / hh}
	ex, ey := self.projmatrix.GetElements()[0], self.projmatrix.GetElements()[4]
	cameraxy := [2]float32{clipxy[0] / ex, clipxy[1] / ey}
	e := self.viewmatrix.GetElements()
	r00, r10, r01, r11, tx, ty := e[0], e[1], e[3], e[4], e[6], e[7]
	x, y := cameraxy[0]-tx, cameraxy[1]-ty          // inverse of (Translation from WORLD to CAMERA)
	wxy := [2]float32{r00*x + r10*y, r01*x + r11*y} // inverse of (Rotation from WORLD to CAMERA)
	return wxy
}

func (self *Camera) UnprojectCanvasDeltaToWorld(deltaxy [2]int) [2]float32 {
	hw, hh := (float32(self.wh[0]) / 2), (float32(self.wh[1]) / 2)
	clip_delta := [2]float32{float32(deltaxy[0]) / hw, -float32(deltaxy[1]) / hh}
	ex, ey := self.projmatrix.GetElements()[0], self.projmatrix.GetElements()[4]
	cam_delta := [2]float32{clip_delta[0] / ex, clip_delta[1] / ey}
	e := self.viewmatrix.GetElements()
	r00, r10, r01, r11 := e[0], e[1], e[3], e[4]
	wx := r00*cam_delta[0] + r10*cam_delta[1] // inverse of (Rotation from WORLD to CAMERA)
	wy := r01*cam_delta[0] + r11*cam_delta[1]
	return [2]float32{wx, wy}
}
