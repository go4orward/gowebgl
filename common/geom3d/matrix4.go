package geom3d

type Matrix4 struct {
	elements [16]float32
}

func NewMatrix4() *Matrix4 {
	mtx := Matrix4{elements: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}} // identity matrix
	return &mtx
}

func (self *Matrix4) GetElements() [16]float32 {
	return self.elements
}

// ----------------------------------------------------------------------------
// Setting element values
// ----------------------------------------------------------------------------

func (self *Matrix4) Set(
	v00 float32, v01 float32, v02 float32, v03 float32,
	v10 float32, v11 float32, v12 float32, v13 float32,
	v20 float32, v21 float32, v22 float32, v23 float32,
	v30 float32, v31 float32, v32 float32, v33 float32) *Matrix4 {
	self.elements = [16]float32{v00, v01, v02, v03, v10, v11, v12, v13, v20, v21, v22, v23, v30, v31, v32, v33}
	return self
}

func (self *Matrix4) SetTranspose() *Matrix4 {
	e := &self.elements // reference
	e[1], e[4] = e[4], e[1]
	e[2], e[8] = e[8], e[2]   // [ 0], [ 1], [ 2], [ 3]
	e[3], e[12] = e[12], e[3] // [ 4], [ 5], [ 6], [ 7]
	e[6], e[9] = e[9], e[6]   // [ 8], [ 9], [10], [11]
	e[7], e[13] = e[13], e[7] // [12], [13], [14], [15]
	e[11], e[14] = e[14], e[11]
	return self
}

func (self *Matrix4) SetCopy(m *Matrix4) *Matrix4 {
	self.elements = m.elements
	return self
}

func (self *Matrix4) SetMultiplyRight(matrix *Matrix4) *Matrix4 {
	o := self.elements    // copy
	m := &matrix.elements // reference
	e := &self.elements
	e[0] = o[0]*m[0] + o[1]*m[4] + o[2]*m[8] + o[3]*m[12]  // [ 0], [ 1], [ 2], [ 3]
	e[1] = o[0]*m[1] + o[1]*m[5] + o[2]*m[9] + o[3]*m[13]  // [ 4], [ 5], [ 6], [ 7]
	e[2] = o[0]*m[2] + o[1]*m[6] + o[2]*m[10] + o[3]*m[14] // [ 8], [ 9], [10], [11]
	e[3] = o[0]*m[3] + o[1]*m[7] + o[2]*m[11] + o[3]*m[15] // [12], [13], [14], [15]
	e[4] = o[4]*m[0] + o[5]*m[4] + o[6]*m[8] + o[7]*m[12]
	e[5] = o[4]*m[1] + o[5]*m[5] + o[6]*m[9] + o[7]*m[13]
	e[6] = o[4]*m[2] + o[5]*m[6] + o[6]*m[10] + o[7]*m[14]
	e[7] = o[4]*m[3] + o[5]*m[7] + o[6]*m[11] + o[7]*m[15]
	e[8] = o[8]*m[0] + o[9]*m[4] + o[10]*m[8] + o[11]*m[12]
	e[9] = o[8]*m[1] + o[9]*m[5] + o[10]*m[9] + o[11]*m[13]
	e[10] = o[8]*m[2] + o[9]*m[6] + o[10]*m[10] + o[11]*m[14]
	e[11] = o[8]*m[3] + o[9]*m[7] + o[10]*m[11] + o[11]*m[15]
	e[12] = o[12]*m[0] + o[13]*m[4] + o[14]*m[8] + o[15]*m[12]
	e[13] = o[12]*m[1] + o[13]*m[5] + o[14]*m[9] + o[15]*m[13]
	e[14] = o[12]*m[2] + o[13]*m[6] + o[14]*m[10] + o[15]*m[14]
	e[15] = o[12]*m[3] + o[13]*m[7] + o[14]*m[11] + o[15]*m[15]
	return self
}

func (self *Matrix4) SetMultiplyMatrices(matrices ...*Matrix4) *Matrix4 {
	for i, m := range matrices {
		if i == 0 {
			self.SetCopy(m)
		} else {
			self.SetMultiplyRight(m)
		}
	}
	return self
}

// ----------------------------------------------------------------------------
// Creating new matrix
// ----------------------------------------------------------------------------

func (self *Matrix4) Transpose() *Matrix4 {
	o := &self.elements // reference
	return &Matrix4{elements: [16]float32{
		o[0], o[4], o[8], o[12],
		o[1], o[5], o[9], o[13],
		o[2], o[6], o[10], o[14],
		o[3], o[7], o[11], o[15]}}
}

func (self *Matrix4) MultiplyRight(matrix *Matrix4) *Matrix4 {
	o := &self.elements   // reference
	m := &matrix.elements // reference
	return &Matrix4{elements: [16]float32{
		o[0]*m[0] + o[1]*m[4] + o[2]*m[8] + o[3]*m[12],  // [ 0], [ 1], [ 2], [ 3]
		o[0]*m[1] + o[1]*m[5] + o[2]*m[9] + o[3]*m[13],  // [ 4], [ 5], [ 6], [ 7]
		o[0]*m[2] + o[1]*m[6] + o[2]*m[10] + o[3]*m[14], // [ 8], [ 9], [10], [11]
		o[0]*m[3] + o[1]*m[7] + o[2]*m[11] + o[3]*m[15], // [12], [13], [14], [15]
		o[4]*m[0] + o[5]*m[4] + o[6]*m[8] + o[7]*m[12],  // 2nd row
		o[4]*m[1] + o[5]*m[5] + o[6]*m[9] + o[7]*m[13],
		o[4]*m[2] + o[5]*m[6] + o[6]*m[10] + o[7]*m[14],
		o[4]*m[3] + o[5]*m[7] + o[6]*m[11] + o[7]*m[15],
		o[8]*m[0] + o[9]*m[4] + o[10]*m[8] + o[11]*m[12], // 3rd
		o[8]*m[1] + o[9]*m[5] + o[10]*m[9] + o[11]*m[13],
		o[8]*m[2] + o[9]*m[6] + o[10]*m[10] + o[11]*m[14],
		o[8]*m[3] + o[9]*m[7] + o[10]*m[11] + o[11]*m[15],
		o[12]*m[0] + o[13]*m[4] + o[14]*m[8] + o[15]*m[12], // 4th
		o[12]*m[1] + o[13]*m[5] + o[14]*m[9] + o[15]*m[13],
		o[12]*m[2] + o[13]*m[6] + o[14]*m[10] + o[15]*m[14],
		o[12]*m[3] + o[13]*m[7] + o[14]*m[11] + o[15]*m[15]}}
}
