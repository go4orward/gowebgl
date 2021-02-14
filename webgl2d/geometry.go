package webgl2d

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

type Geometry struct {
	verts [][2]float32
	edges [][]uint32
	faces [][]uint32
	tuvs  [][]float32 // texture uv coordinates (PER_FACE [nfaces][6] or PER_VERT [nverts][2])

	data_buffer_vpoints []float32 // serialized data buffer for vertex points : COORD[]
	data_buffer_fpoints []float32 // serialized data buffer for PER_FACE vertex points : COORD[2] + (UV[2])
	data_buffer_lines   []uint32  // serialized data buffer for edge lines : vidx[2]
	data_buffer_faces   []uint32  // serialized data buffer for face triangles : vidx[3]

	// Note that, for PER_FACE texture UV-coordinates, vertices are duplicated for each face
	fpoint_vidx_list  []uint32 // index of vertex_list of each face after PER_FACE data duplication
	fpoint_vert_total int      // total count of vertices after PER_FACE data duplication
	fpoint_info       [3]int   // data size of a point (for triangles) : [ stride, xyz_offset, uv_offset ]
	vpoint_info       [3]int   // data size of a point (for points & lines)

	webgl_buffer_vpoints js.Value // WebGL data buffer for data_buffer_vpoints (points for vertices)
	webgl_buffer_fpoints js.Value // WebGL data buffer for data_buffer_fpoints (points for PER_FACE vertices)
	webgl_buffer_lines   js.Value // WebGL data buffer for data_buffer_lines (indices for lines)
	webgl_buffer_faces   js.Value // WebGL data buffer for data_buffer_faces (indices for triangles)
}

func NewGeometry() *Geometry {
	var geometry Geometry
	geometry.Clear(true, true, true)
	return &geometry
}

func (self *Geometry) Clear(geom bool, data_buf bool, webgl_buf bool) *Geometry {
	if geom {
		self.verts = [][2]float32{}
		self.edges = [][]uint32{}
		self.faces = [][]uint32{}
		self.tuvs = [][]float32{}
	}
	if data_buf || geom {
		self.data_buffer_vpoints = nil
		self.data_buffer_fpoints = nil
		self.data_buffer_lines = nil
		self.data_buffer_faces = nil
		self.fpoint_vidx_list = nil
		self.fpoint_vert_total = 0
		self.fpoint_info = [3]int{0, 0, 0}
		self.vpoint_info = [3]int{0, 0, 0}
	}
	if webgl_buf || data_buf || geom {
		self.webgl_buffer_vpoints = js.Null()
		self.webgl_buffer_fpoints = js.Null()
		self.webgl_buffer_lines = js.Null()
		self.webgl_buffer_faces = js.Null()
	}
	return self
}

func (self *Geometry) ShowInfo() {
	wblen := func(b js.Value) string {
		if b.IsNull() {
			return "NULL"
		} else if b.IsUndefined() {
			return "UDEF"
		} else {
			return fmt.Sprintf("%4d", b.Length())
		}
	}
	fmt.Printf("Geometry with %d verts %d edges %d faces\n", len(self.verts), len(self.edges), len(self.faces))
	if len(self.tuvs) > 0 {
		fmt.Printf("  texture UV coords   : [%d][]float32\n", len(self.tuvs))
	}
	fmt.Printf("  data_buffer_vpoints : %4d (WebGL:%s)  psize=%v\n", len(self.data_buffer_vpoints), wblen(self.webgl_buffer_vpoints), self.vpoint_info)
	fmt.Printf("  data_buffer_fpoints : %4d (WebGL:%s)  psize=%v\n", len(self.data_buffer_fpoints), wblen(self.webgl_buffer_fpoints), self.fpoint_info)
	fmt.Printf("  data_buffer_lines   : %4d (WebGL:%s)\n", len(self.data_buffer_lines), wblen(self.webgl_buffer_lines))
	fmt.Printf("  data_buffer_faces   : %4d (WebGL:%s)\n", len(self.data_buffer_faces), wblen(self.webgl_buffer_faces))
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *Geometry) SetVertices(vertices [][2]float32) *Geometry {
	self.verts = vertices
	return self
}

func (self *Geometry) SetEdges(edges [][]uint32) *Geometry {
	self.edges = edges
	return self
}

func (self *Geometry) SetFaces(faces [][]uint32) *Geometry {
	self.faces = faces
	return self
}

func (self *Geometry) AddVertex(coords [2]float32) uint32 {
	vidx := len(self.verts)
	self.verts = append(self.verts, coords)
	return uint32(vidx)
}

func (self *Geometry) AddEdge(edge []uint32) uint32 {
	eidx := len(self.edges)
	self.edges = append(self.edges, edge)
	return uint32(eidx)
}

func (self *Geometry) AddFace(face []uint32) uint32 {
	fidx := len(self.faces)
	self.faces = append(self.faces, face)
	return uint32(fidx)
}

// ----------------------------------------------------------------------------
// Transformation of Vertex Coordinates
// ----------------------------------------------------------------------------

func (self *Geometry) Translate(tx float32, ty float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] += tx
		self.verts[i][1] += ty
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Rotate(angle_in_degree float32) *Geometry {
	rad := float64(angle_in_degree * (math.Pi / 180))
	sin, cos := math.Sin(rad), math.Cos(rad)
	for i := 0; i < len(self.verts); i++ {
		x, y := float64(self.verts[i][0]), float64(self.verts[i][1])
		self.verts[i][0] = float32(cos*x - sin*y)
		self.verts[i][1] = float32(sin*x + cos*y)
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Scale(scale float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] *= scale
		self.verts[i][1] *= scale
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) AppyMatrix(matrix *geom2d.Matrix3) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i] = matrix.MultiplyVector2(self.verts[i])
	}
	self.Clear(false, true, true)
	return self
}

// ----------------------------------------------------------------------------
// Texture UV coordinates
// ----------------------------------------------------------------------------

func (self *Geometry) HasTextureFor(mode string) bool {
	switch mode {
	case "VERTEX":
		return len(self.tuvs) > 0 && len(self.tuvs[0]) == 2
	case "FACE":
		return len(self.tuvs) > 0 && len(self.tuvs[0]) >= 6
	default:
		return self.HasTextureFor("VERTEX") || self.HasTextureFor("FACE")
	}
}

func (self *Geometry) AddTextureUV(tuv []float32) *Geometry {
	if len(tuv) == 0 || len(tuv)%2 == 0 {
		fmt.Printf("Invalid texture coordinates to add : %v\n", tuv)
		return self
	}
	self.tuvs = append(self.tuvs, tuv)
	return self
}

func (self *Geometry) SetTextureUVs(tuvs [][]float32) *Geometry {
	self.tuvs = tuvs
	return self
}

// func (self *Geometry) AllocTextureUVsPerVertex() {
// 	self.tuvs = make([][]float32, len(self.verts)) // self.tuvs == [nverts][2]float32
// 	for i, _ := range self.verts {
// 		self.tuvs[i] = make([]float32, 2)
// 	}
// }

// func (self *Geometry) AllocTextureUVsPerFace() {
// 	self.tuvs = make([][]float32, len(self.faces)) // self.tuvs == [nfaces][fvcount*2]float32
// 	for i, face := range self.faces {
// 		self.tuvs[i] = make([]float32, 2*len(face))
// 	}
// }

// ----------------------------------------------------------------------------
// Triangulation
// ----------------------------------------------------------------------------

func (self *Geometry) get_reverse(face_vlist []uint32) []uint32 {
	new_vlist := make([]uint32, len(face_vlist))
	for i := len(face_vlist) - 1; i >= 0; i-- {
		new_vlist[i] = face_vlist[i]
	}
	return new_vlist
}

func (self *Geometry) splice_indices(a []uint32, pos int, delete_count int, new_entries ...uint32) []uint32 {
	head := a[0:pos]
	tail := a[pos+delete_count:]
	return append(append(head, new_entries...), tail...)
}

func (self *Geometry) get_triangulation(face []uint32) [][]uint32 {
	vindices := make([]uint32, len(face))
	copy(vindices, face)
	newfaces := make([][]uint32, 0)
	vidx, vcount := 0, len(vindices)
	for vcount > 3 {
		i0, i1, i2 := vidx, (vidx+1)%vcount, (vidx+2)%vcount
		v0, v1, v2 := self.verts[vindices[i0]], self.verts[vindices[i1]], self.verts[vindices[i2]]
		if geom2d.IsCCW(v0, v1, v2) {
			point_inside := false
			for j := 0; j < vcount; j++ {
				if j != i0 && j != i1 && j != i2 && geom2d.IsPointInside(self.verts[vindices[j]], v0, v1, v2) {
					point_inside = true
					break
				}
			}
			if !point_inside {
				newfaces = append(newfaces, []uint32{vindices[i0], vindices[i1], vindices[i2]})
				vindices = self.splice_indices(vindices, i1, 1)
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

func (self *Geometry) IsDataBufferReady() bool {
	return len(self.data_buffer_vpoints) > 0 || len(self.data_buffer_fpoints) > 0
}

func (self *Geometry) count_fpoint_vidx_list() int {
	self.fpoint_vert_total = 0
	self.fpoint_vidx_list = make([]uint32, len(self.faces))
	for i := 0; i < len(self.faces); i++ {
		self.fpoint_vidx_list[i] = uint32(self.fpoint_vert_total)
		self.fpoint_vert_total += len(self.faces[i])
	}
	return self.fpoint_vert_total
}

func (self *Geometry) get_fpoint_new_vidx(fidx int, i int) int {
	if self.fpoint_vidx_list == nil {
		self.count_fpoint_vidx_list()
	}
	return int(self.fpoint_vidx_list[fidx]) + i
}

func (self *Geometry) copy_buffer(buf []float32, new_vidx int, vidx int, tuv_size int, tuv_idx int) {
	pos := new_vidx * (2 + tuv_size)
	if true {
		buf[pos+0], buf[pos+1] = self.verts[vidx][0], self.verts[vidx][1]
		pos += 2
	}
	if tuv_size == 2 {
		buf[pos+0], buf[pos+1] = self.tuvs[tuv_idx][0], self.tuvs[tuv_idx][1]
		pos += 2
	}
}

func (self *Geometry) BuildDataBuffers(for_points bool, for_lines bool, for_faces bool) {
	// create data buffer for vertex points
	self.data_buffer_vpoints, self.vpoint_info = nil, [3]int{0, 0, 0}
	self.data_buffer_fpoints, self.fpoint_info = nil, [3]int{0, 0, 0}
	points_per_face := false
	if for_faces {
		if self.HasTextureFor("FACE") {
			points_per_face = true
			self.count_fpoint_vidx_list()
			self.fpoint_info = [3]int{(2 + 2), (0), (3)} // size, xyz_offset, uv_offset
			self.data_buffer_fpoints = make([]float32, self.fpoint_vert_total*self.fpoint_info[0])
			for fidx, face := range self.faces {
				for i := 0; i < len(face); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face[i])
					self.copy_buffer(self.data_buffer_fpoints, new_vidx, vidx, 2, fidx)
				}
			}
		} else if self.HasTextureFor("VERTEX") {
			self.fpoint_info = [3]int{(2 + 2), (0), (3)} // size, xyz_offset, uv_offset
			self.data_buffer_fpoints = make([]float32, len(self.verts)*self.fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.copy_buffer(self.data_buffer_fpoints, vidx, vidx, 2, vidx)
			}
			self.data_buffer_vpoints = self.data_buffer_fpoints
			self.vpoint_info = self.fpoint_info
		} else {
			self.fpoint_info = [3]int{(2 + 0), (0), (0)} // size, xyz_offset, uv_offset
			self.data_buffer_fpoints = make([]float32, len(self.verts)*self.fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.copy_buffer(self.data_buffer_fpoints, vidx, vidx, 0, 0)
			}
			self.data_buffer_vpoints = self.data_buffer_fpoints
			self.vpoint_info = self.fpoint_info
		}
	} else {
		self.data_buffer_vpoints = nil
	}
	if (for_points || for_lines) && self.data_buffer_vpoints == nil {
		self.vpoint_info = [3]int{2, 0, 0}
		self.data_buffer_vpoints = make([]float32, len(self.verts)*self.vpoint_info[0])
		for vidx := 0; vidx < len(self.verts); vidx++ {
			self.copy_buffer(self.data_buffer_vpoints, vidx, vidx, 0, 0)
		}
	}
	// if self.fpoint_info[0] == 0 {
	// 	self.fpoint_info = [3]int{self.vpoint_info, 0, 0}
	// }
	// create data buffer for line drawings
	if for_lines {
		segment_count := 0
		for _, edge := range self.edges {
			segment_count += len(edge) - 1
		}
		self.data_buffer_lines = make([]uint32, segment_count*2)
		epos := 0
		for _, edge := range self.edges {
			for i := 1; i < len(edge); i++ {
				self.data_buffer_lines[epos+0] = edge[i-1]
				self.data_buffer_lines[epos+1] = edge[i]
				epos += 2
			}
		}
	} else {
		self.data_buffer_lines = nil
	}
	// create data buffer for surface drawings
	if for_faces {
		triangle_count := 0
		for _, face := range self.faces {
			triangle_count += len(face) - 2
		}
		find_index_in_face := func(vidx uint32, face []uint32) int {
			for i := 0; i < len(face); i++ {
				if vidx == face[i] {
					return i
				}
			}
			return 0
		}
		self.data_buffer_faces = make([]uint32, triangle_count*3)
		tpos := 0
		for fidx, face := range self.faces { // []vidx
			triangles := self.get_triangulation(face) // [][3]vidx
			for _, triangle := range triangles {      // [3]vidx
				if points_per_face { // vertex index has been changed due to PER_FACE duplication
					vidx_stt := self.get_fpoint_new_vidx(fidx, 0)
					self.data_buffer_faces[tpos+0] = uint32(vidx_stt + find_index_in_face(triangle[0], face))
					self.data_buffer_faces[tpos+1] = uint32(vidx_stt + find_index_in_face(triangle[1], face))
					self.data_buffer_faces[tpos+2] = uint32(vidx_stt + find_index_in_face(triangle[2], face))
				} else { // vertex index was preserved
					self.data_buffer_faces[tpos+0] = triangle[0]
					self.data_buffer_faces[tpos+1] = triangle[1]
					self.data_buffer_faces[tpos+2] = triangle[2]
				}
				tpos += 3
			}
		}
	} else {
		self.data_buffer_faces = nil
	}
	self.Clear(false, false, true)
}

func (self *Geometry) BuildDataBuffersForWireframe() {
	if self.data_buffer_vpoints == nil {
		// create data buffer for vertex points, only if necessary
		self.data_buffer_vpoints = make([]float32, len(self.verts)*2)
		vpos := 0
		for _, xy := range self.verts {
			self.data_buffer_vpoints[vpos+0] = xy[0]
			self.data_buffer_vpoints[vpos+1] = xy[1]
			vpos += 2
		}
		self.vpoint_info = [3]int{2, 0, 0}
	}
	// create data buffer for edges, by extracting wireframe from faces
	self.data_buffer_lines = make([]uint32, 0)
	for _, face := range self.faces {
		triangles := self.get_triangulation(face)
		for _, t := range triangles {
			self.data_buffer_lines = append(self.data_buffer_lines, t[0], t[1], t[1], t[2], t[2], t[0])
		}
	}
	self.Clear(false, false, true)
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *Geometry) IsWebGLBufferReady() bool {
	return !self.webgl_buffer_vpoints.IsNull()
}

func (self *Geometry) build_webgl_buffers(wctx *common.WebGLContext, for_points bool, for_lines bool, for_faces bool) {
	// THIS FUCNTION IS MEANT TO BE CALLED BY RENDERER. NO NEED TO BE EXPORTED
	context := wctx.GetContext()     // js.Value
	constants := wctx.GetConstants() // *common.Constants
	if for_points && self.data_buffer_vpoints != nil {
		self.webgl_buffer_vpoints = context.Call("createBuffer", constants.ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, self.webgl_buffer_vpoints)
		var vertices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_vpoints)
		context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_vpoints = js.Null()
	}
	if for_lines && self.data_buffer_lines != nil {
		self.webgl_buffer_lines = context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, self.webgl_buffer_lines)
		var indices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_lines)
		context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_lines = js.Null()
	}
	if for_faces && self.data_buffer_faces != nil {
		if self.data_buffer_fpoints != nil {
			self.webgl_buffer_fpoints = context.Call("createBuffer", constants.ARRAY_BUFFER)
			context.Call("bindBuffer", constants.ARRAY_BUFFER, self.webgl_buffer_fpoints)
			var points_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_fpoints)
			context.Call("bufferData", constants.ARRAY_BUFFER, points_array, constants.STATIC_DRAW)
			context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
		}
		self.webgl_buffer_faces = context.Call("createBuffer", constants.ELEMENT_ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, self.webgl_buffer_faces)
		var indices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer_faces)
		context.Call("bufferData", constants.ELEMENT_ARRAY_BUFFER, indices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ELEMENT_ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer_faces = js.Null()
	}
}

func (self *Geometry) GetWebGLBuffer(mode string) (js.Value, int, [3]int) {
	switch mode {
	case "POINTS":
		if self.data_buffer_fpoints == nil {
			return self.webgl_buffer_vpoints, len(self.data_buffer_vpoints), self.vpoint_info
		} else if self.data_buffer_vpoints == nil {
			return self.webgl_buffer_fpoints, len(self.data_buffer_fpoints), self.fpoint_info
		} else {
			return self.webgl_buffer_fpoints, len(self.data_buffer_fpoints), self.fpoint_info
		}
	case "LINES":
		return self.webgl_buffer_lines, len(self.data_buffer_lines), [3]int{0, 0, 0}
	case "TRIANGLES":
		return self.webgl_buffer_faces, len(self.data_buffer_faces), [3]int{0, 0, 0}
	default:
		fmt.Printf("Invalid mode '%s' for GetWebGLBuffer()\n")
		return js.Null(), 0, [3]int{0, 0, 0}
	}
}
