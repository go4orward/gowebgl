package webgl2d

import (
	"math"
)

func (self *Geometry) LoadTriangle(size float32) *Geometry {
	self.LoadPolygon(3, size, -30)
	return self
}

func (self *Geometry) LoadRectangle(size float32) *Geometry {
	self.Clear(true, true, true)
	self.AddVertices([][2]float32{{0, 0}, {size, 0}, {size, size}, {0, size}})
	self.AddFaces([][]uint32{{0, 1, 2}, {0, 2, 3}})
	self.Translate(-size/2, -size/2)
	return self
}

func (self *Geometry) LoadPolygon(n int, radius float32, starting_angle_in_degree float32) *Geometry {
	self.Clear(true, true, true)
	radian := float64(starting_angle_in_degree * (math.Pi / 180.0))
	radian_step := (2 * math.Pi) / float64(n)
	face_indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		self.AddVertex([2]float32{radius * float32(math.Cos(radian)), radius * float32(math.Sin(radian))})
		face_indices[i] = uint32(i)
		// fmt.Printf("angle: %4.0f %v\n", radian*(180.0/math.Pi), self.verts)
		radian += radian_step
	}
	self.AddFace(face_indices)
	return self
}
