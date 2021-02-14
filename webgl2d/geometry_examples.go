package webgl2d

import (
	"math"
)

func NewGeometry_Triangle(size float32) *Geometry {
	return NewGeometry_Polygon(3, size, -30)
}

func NewGeometry_Rectangle(size float32) *Geometry {
	geometry := NewGeometry() // create an empty geometry
	geometry.SetVertices([][2]float32{{0, 0}, {size, 0}, {size, size}, {0, size}})
	geometry.SetFaces([][]uint32{{0, 1, 2}, {0, 2, 3}})
	geometry.Translate(-size/2, -size/2)
	return geometry
}

func NewGeometry_Polygon(n int, radius float32, starting_angle_in_degree float32) *Geometry {
	geometry := NewGeometry() // create an empty geometry
	radian := float64(starting_angle_in_degree * (math.Pi / 180.0))
	radian_step := (2 * math.Pi) / float64(n)
	face_indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		geometry.AddVertex([2]float32{radius * float32(math.Cos(radian)), radius * float32(math.Sin(radian))})
		face_indices[i] = uint32(i)
		radian += radian_step
	}
	geometry.AddFace(face_indices)
	return geometry
}
