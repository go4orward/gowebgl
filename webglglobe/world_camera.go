package webglglobe

import (
	"math"

	"github.com/go4orward/gowebgl/geom3d"
	"github.com/go4orward/gowebgl/webgl3d"
)

type WorldCamera struct {
	gcam *webgl3d.Camera // Perspective Globe Camera
	// TODO: Add support for Orthographic projection camera
	// TODO: Add more cameras for flat world map projections
}

func NewWorldCamera(wh [2]int, fov float32, zoom float32) *WorldCamera {
	camera := webgl3d.NewPerspectiveCamera(wh, fov, zoom) // AspectRatio, FieldOfView, ZoomLevel
	self := WorldCamera{gcam: camera}
	return &self
}

func (self *WorldCamera) ShowInfo() {
	self.gcam.ShowInfo()
}

// ----------------------------------------------------------------------------
// Camera Internal Parameters
// ----------------------------------------------------------------------------

func (self *WorldCamera) SetAspectRatio(width int, height int) *WorldCamera {
	self.gcam.SetAspectRatio(width, height)
	return self
}

func (self *WorldCamera) SetZoom(zoom float32) *WorldCamera {
	self.gcam.SetZoom(zoom)
	return self
}

// ----------------------------------------------------------------------------
// Camera Pose
// ----------------------------------------------------------------------------

func (self *WorldCamera) SetPoseByLonLat(lon float32, lat float32, dist float32) *WorldCamera {
	Twc := GetXYZFromLonLat(lon, lat, dist)              // Camera center in WORLD space
	coslon := float32(math.Cos(float64(lon) * InRadian)) // cos(λ)
	sinlon := float32(math.Sin(float64(lon) * InRadian)) // sin(λ)
	camX := [3]float32{-sinlon, +coslon, 0}              // this prevents UP vector singularity at poles
	camZ := geom3d.Normalize(Twc)                        // camera's Z axis points backward (away from view frustum)
	camY := geom3d.CrossAB(camZ, camX)
	self.gcam.SetPoseWithCameraAxes(camX, camY, camZ, Twc)
	return self
}

func (self *WorldCamera) RotateAroundGlobe(horizontal_angle float32, vertical_angle float32) *WorldCamera {
	self.gcam.RotateAroundPoint(10, horizontal_angle, vertical_angle)
	self.RotateByRollToHeadUpNorth()
	return self
}

func (self *WorldCamera) RotateByRollToHeadUpNorth() *WorldCamera {
	e := self.gcam.GetViewMatrix().GetElements()
	// Note that {e[8],e[9]} is the NORTH (0,0,1) projected onto XY plane of Camera axes in WORLD space.
	if e[8]*e[8]+e[9]*e[9] > 0.01 { // Now compare {e[8],e[9]} with Y direction (90°) in CAMERA space.
		roll := 90 - float32(math.Atan2(float64(e[9]), float64(e[8])))*InDegree
		self.gcam.RotateByRoll(roll)
	}
	return self
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------
