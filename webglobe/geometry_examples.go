package webglobe

import (
	"math"

	"github.com/go4orward/gowebgl/webgl3d"
)

const InRadian = (math.Pi / 180.0)
const InDegree = (180.0 / math.Pi)

func NewGeometry_Globe(radius float32, wsegs int, hsegs int) *webgl3d.Geometry {
	// Globe (sphere) geometry with UV coordinates per vertex (to be used with a texture image)
	//   Note that multiple vertices are assigned to north/south poles, as well as 0/360 longitude.
	//   This approach results in more efficient data buffers than a simple sphere,
	//   since we can build the buffers with single point per vertex, without any repetition.
	geometry := webgl3d.NewGeometry()
	wnum, hnum := wsegs+1, hsegs+1
	wstep := math.Pi * 2.0 / float32(wsegs)
	hstep := math.Pi / float32(hsegs)
	for i := 0; i < wnum; i++ {
		lon := wstep * float32(i) // longitude (λ)
		for j := 0; j < hnum; j++ {
			lat := -math.Pi/2.0 + hstep*float32(j) // latitude (φ)
			xyz := GetXYZFromLL(lon, lat, radius)
			geometry.AddVertex(xyz)
			if pole := (j == 0 || j == hsegs); pole {
				geometry.AddTextureUV([]float32{(float32(i) + 0.5) / float32(wsegs), 1.0 - float32(j)/float32(hsegs)})
			} else {
				geometry.AddTextureUV([]float32{float32(i) / float32(wsegs), 1.0 - float32(j)/float32(hsegs)})
			}
		}
	}
	for i := 0; i < wnum-1; i++ { // quadratic faces on the side (triangles for south/north poles)
		for j := 0; j < hnum-1; j++ {
			start := uint32((i+0)*hnum + j)
			wnext := uint32((i+1)*hnum + j)
			if spole := (j == 0); spole {
				geometry.AddFace([]uint32{start, wnext + 1, start + 1})
			} else if npole := (j == hsegs-1); npole {
				geometry.AddFace([]uint32{start, wnext + 0, start + 1})
			} else {
				geometry.AddFace([]uint32{start, wnext, wnext + 1, start + 1})
			}
			// fmt.Printf("adding Face [ %d %d %d %d ]\n", start, wnext, wnext+1, start+1)
		}
	}
	return geometry
}
