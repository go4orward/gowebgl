package geom2d

import "math"

func AddAB(a [2]float32, b [2]float32) [2]float32 {
	return [2]float32{a[0] + b[0], a[1] + b[1]}
}

func SubAB(a [2]float32, b [2]float32) [2]float32 {
	return [2]float32{a[0] - b[0], a[1] - b[1]}
}

func CrossAB(a [2]float32, b [2]float32) float32 {
	return a[0]*b[1] - a[1]*b[0] // in 2D, (ax,ay,0) x (bx,by,0) = (0,0,ax*by-ay*bx)
}

func Length(v [2]float32) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1])))
}

func Normalize(v [2]float32) [2]float32 {
	len := Length(v)
	return [2]float32{v[0] / len, v[1] / len}
}

func IsCCW(v0 [2]float32, v1 [2]float32, v2 [2]float32) bool {
	v01 := SubAB(v1, v0)
	v02 := SubAB(v2, v0)
	return CrossAB(v01, v02) > 0
}

func IsPointInside(p [2]float32, v0 [2]float32, v1 [2]float32, v2 [2]float32) bool {
	p0, p1, p2 := SubAB(v0, p), SubAB(v1, p), SubAB(v2, p)
	c01, c12, c13 := CrossAB(p0, p1), CrossAB(p1, p2), CrossAB(p2, p0)
	return (c01 > 0 && c12 > 0 && c13 > 0) || (c01 < 0 && c12 < 0 && c13 < 0)
}

// public static isPointInTriangle(point: number[], v0: number[], v1: number[], v2: number[], strictly_inside: boolean = false): boolean {
// 	let p0 = Point.subAB(v0, point), p1 = Point.subAB(v1, point), p2 = Point.subAB(v2, point);
// 	let c01 = Point.crossAB(p0, p1), c12 = Point.crossAB(p1, p2), c20 = Point.crossAB(p2, p0);
// 	let d012 = Point.dotAB(c01,c12), d120 = Point.dotAB(c12, c20), d201 = Point.dotAB(c20, c01);
// 	if (d012 > 0 && d120 > 0 && d201 > 0) return true;  // point is strictly inside the triangle
// 	if (strictly_inside || d012 * d120 * d201 != 0) return false;    // point is not on any side
// 	if (Point.isZero(c01) && d120 < 0) return false;    // point is on side 01, but it's outside
// 	if (Point.isZero(c12) && d201 < 0) return false;    // point is on side 12, but it's outside
// 	if (Point.isZero(c20) && d012 < 0) return false;    // point is on side 20, but it's outside
// 	return true;        // point is on the border, and it's not outside
// }
