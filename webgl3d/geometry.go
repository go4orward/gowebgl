package webgl3d

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom3d"
)

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

type Geometry struct {
	verts [][3]float32
	edges [][]uint32
	faces [][]uint32

	data_buffer_vert  []float32 // serialized data buffer for vertices
	data_buffer_edge  []uint32  // serialized data buffer for edges
	data_buffer_face  []uint32  // serialized data buffer for faces
	webgl_buffer_vert js.Value  // WebGL data buffers for vertices
	webgl_buffer_edge js.Value  // WebGL data buffers for edges
	webgl_buffer_face js.Value  // WebGL data buffers for faces
}

func NewGeometry() *Geometry {
	var geometry Geometry
	geometry.Clear(true, true, true)
	return &geometry
}

func (self *Geometry) Clear(geom bool, data_buf bool, webgl_buf bool) *Geometry {
	if geom {
		self.verts = make([][3]float32, 0)
		self.edges = make([][]uint32, 0)
		self.faces = make([][]uint32, 0)
	}
	if data_buf || geom {
		self.data_buffer_vert = nil
		self.data_buffer_edge = nil
		self.data_buffer_face = nil
	}
	if webgl_buf || data_buf || geom {
		self.webgl_buffer_vert = js.Null()
		self.webgl_buffer_edge = js.Null()
		self.webgl_buffer_face = js.Null()
	}
	return self
}

func (self *Geometry) AddVertex(coords [3]float32) *Geometry {
	self.verts = append(self.verts, coords)
	return self
}

func (self *Geometry) AddVertices(coords [][3]float32) *Geometry {
	self.verts = coords
	return self
}

func (self *Geometry) AddEdge(edge []uint32) *Geometry {
	self.edges = append(self.edges, edge)
	return self
}

func (self *Geometry) AddEdges(indices [][]uint32) *Geometry {
	self.edges = indices
	return self
}

func (self *Geometry) AddFace(face []uint32) *Geometry {
	self.faces = append(self.faces, face)
	return self
}

func (self *Geometry) AddFaces(indices [][]uint32) *Geometry {
	self.faces = indices
	return self
}

func (self *Geometry) Translate(tx float32, ty float32, tz float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] += tx
		self.verts[i][1] += ty
		self.verts[i][2] += tz
	}
	return self
}

func (self *Geometry) Rotate(axis []float32, angle_in_degree float32) *Geometry {
	// rad := float64(angle_in_degree * (math.Pi / 180))
	// sin, cos := math.Sin(rad), math.Cos(rad)
	// for i := 0; i < len(self.verts); i++ {
	// 	x, y := float64(self.verts[i][0]), float64(self.verts[i][1])
	// 	self.verts[i][0] = float32(cos*x - sin*y)
	// 	self.verts[i][1] = float32(sin*x + cos*y)
	// }
	return self
}

func (self *Geometry) Scale(sx float32, sy float32, sz float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] *= sx
		self.verts[i][1] *= sy
		self.verts[i][2] *= sz
	}
	return self
}

func (self *Geometry) AppyMatrix(matrix *geom3d.Matrix4) *Geometry {
	// for i := 0; i < len(self.verts); i++ {
	// 	self.verts[i][0] *= scale
	// 	self.verts[i][1] *= scale
	// }
	return self
}

func (self *Geometry) ShowInfo() {
	fmt.Printf("Geometry with %d verts %d edges %d faces\n", len(self.verts), len(self.edges), len(self.faces))
	if self.data_buffer_vert != nil {
		fmt.Printf("  Data  Buffers : %d %d %d\n", len(self.data_buffer_vert), len(self.data_buffer_edge), len(self.data_buffer_face))
	}
	if !self.webgl_buffer_vert.IsNull() {
		fmt.Printf("  WebGL Buffers : %d %d %d\n", len(self.data_buffer_vert), len(self.data_buffer_edge), len(self.data_buffer_face))
	}
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (self *Geometry) TriangulateFace(face []uint32) [][]uint32 {
	vindices := make([]uint32, len(face))
	copy(vindices, face)
	newfaces := make([][]uint32, 0)
	vidx, vcount := 0, len(vindices)
	for vcount > 3 {
		i0, i1, i2 := vidx, (vidx+1)%vcount, (vidx+2)%vcount
		v0, v1, v2 := self.verts[vindices[i0]], self.verts[vindices[i1]], self.verts[vindices[i2]]
		if geom3d.IsCCW(v0, v1, v2) {
			point_inside := false
			for j := 0; j < vcount; j++ {
				if j != i0 && j != i1 && j != i2 && geom3d.IsPointInside(self.verts[vindices[j]], v0, v1, v2) {
					point_inside = true
					break
				}
			}
			if !point_inside {
				newfaces = append(newfaces, []uint32{vindices[i0], vindices[i1], vindices[i2]})
				vindices = geom3d.SpliceUint32(vindices, i1, 1)
			}
		}
		vcount = len(vindices)
		vidx = (vidx + 1) % vcount
	}
	newfaces = append(newfaces, vindices)
	return newfaces
}

// ----------------------------------------------------------------------------
// Build Data Buffers (serialized)
// ----------------------------------------------------------------------------

func (self *Geometry) BuildDataBuffers(for_verts bool, for_edges bool, for_faces bool) {
	// create data buffer for vertex points
	if for_verts {
		self.data_buffer_vert = make([]float32, len(self.verts)*3)
		vpos := 0
		for _, xyz := range self.verts {
			self.data_buffer_vert[vpos+0] = xyz[0]
			self.data_buffer_vert[vpos+1] = xyz[1]
			self.data_buffer_vert[vpos+2] = xyz[2]
			vpos += 3
		}
	} else {
		self.data_buffer_vert = nil
	}
	// create data buffer for edges
	if for_edges {
		segment_count := 0
		for _, edge := range self.edges {
			segment_count += len(edge) - 1
		}
		self.data_buffer_edge = make([]uint32, segment_count*2)
		epos := 0
		for _, edge := range self.edges {
			for i := 1; i < len(edge); i++ {
				self.data_buffer_edge[epos+0] = edge[i-1]
				self.data_buffer_edge[epos+0] = edge[i]
				epos += 2
			}
		}
	} else {
		self.data_buffer_edge = nil
	}
	// create data buffer for faces
	if for_faces {
		triangle_count := 0
		for _, face := range self.faces {
			triangle_count += len(face) - 2
		}
		self.data_buffer_face = make([]uint32, triangle_count*3)
		tpos := 0
		for _, face := range self.faces {
			triangles := self.TriangulateFace(face)
			for _, triangle := range triangles {
				self.data_buffer_face[tpos+0] = triangle[0]
				self.data_buffer_face[tpos+1] = triangle[1]
				self.data_buffer_face[tpos+2] = triangle[2]
				tpos += 3
			}
		}
	} else {
		self.data_buffer_face = nil
	}
	self.Clear(false, false, true)
}

func (self *Geometry) BuildDataBuffersForWireframe() {
	if self.data_buffer_vert == nil {
		// create data buffer for vertex points, only if necessary
		self.data_buffer_vert = make([]float32, len(self.verts)*3)
		vpos := 0
		for _, xyz := range self.verts {
			self.data_buffer_vert[vpos+0] = xyz[0]
			self.data_buffer_vert[vpos+1] = xyz[1]
			self.data_buffer_vert[vpos+2] = xyz[2]
			vpos += 3
		}
	}
	// create data buffer for edges, by extracting wireframe from faces
	self.data_buffer_edge = make([]uint32, 0)
	for _, face := range self.faces {
		triangles := self.TriangulateFace(face)
		for _, triangle := range triangles {
			self.data_buffer_edge = append(self.data_buffer_edge, triangle[0], triangle[1])
			self.data_buffer_edge = append(self.data_buffer_edge, triangle[1], triangle[2])
			self.data_buffer_edge = append(self.data_buffer_edge, triangle[2], triangle[0])
		}
	}
	self.Clear(false, false, true)
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *Geometry) IsWebGLBufferReady() bool {
	return !self.webgl_buffer_vert.IsNull()
}

func (self *Geometry) build_webgl_buffers(wctx *common.WebGLContext, for_vert bool, for_edge bool, for_face bool) {
	context := wctx.GetContext()     // js.Value
	constants := wctx.GetConstants() // *common.Constants
	if for_vert && self.data_buffer_vert != nil {
		self.webgl_buffer_vert = context.Call("createBuffer", constants.ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, self.webgl_buffer_vert)
		var vertices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_vert)
		context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_vert = js.Null()
	}
	if for_edge && self.data_buffer_edge != nil {
		self.webgl_buffer_edge = context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, self.webgl_buffer_edge)
		var indices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_edge)
		context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_edge = js.Null()
	}
	if for_face && self.data_buffer_face != nil {
		self.webgl_buffer_face = context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, self.webgl_buffer_face)
		var indices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_face)
		context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_face = js.Null()
	}
}

func (self *Geometry) GetWebGLBufferToDraw(mode string) js.Value {
	switch mode {
	case "POINTS":
		return self.webgl_buffer_vert
	case "LINES":
		return self.webgl_buffer_edge
	case "TRIANGLES":
		return self.webgl_buffer_face
	default:
		return self.webgl_buffer_face
	}
}

func (self *Geometry) GetWebGLCountToDraw(mode string) int {
	switch mode {
	case "POINTS":
		return len(self.data_buffer_vert)
	case "LINES":
		return len(self.data_buffer_edge)
	case "TRIANGLES":
		return len(self.data_buffer_face)
	default:
		return len(self.data_buffer_face)
	}
}
