package geom2d

type Matrix3 struct {
	elements [9]float32
}

func NewMatrix3() *Matrix3 {
	mtx := Matrix3{elements: [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}} // identity matrix
	return &mtx
}

func (self *Matrix3) GetElements() [9]float32 {
	return self.elements
}

// ----------------------------------------------------------------------------
// Setting element values
// ----------------------------------------------------------------------------

func (self *Matrix3) Set(v00 float32, v01 float32, v02 float32, v10 float32, v11 float32, v12 float32, v20 float32, v21 float32, v22 float32) *Matrix3 {
	self.elements = [9]float32{v00, v01, v02, v10, v11, v12, v20, v21, v22}
	return self
}

func (self *Matrix3) SetCopy(m *Matrix3) *Matrix3 {
	self.elements = m.elements // copy values
	return self
}

func (self *Matrix3) SetTranspose() *Matrix3 {
	e := &self.elements     // reference
	e[1], e[3] = e[3], e[1] // [0], [1], [2]
	e[2], e[6] = e[6], e[2] // [3], [4], [5]
	e[5], e[7] = e[7], e[5] // [6], [7], [8]
	return self
}

func (self *Matrix3) SetMultiplyRight(matrix *Matrix3) *Matrix3 {
	o := self.elements    // copy
	m := &matrix.elements // reference
	self.elements[0] = o[0]*m[0] + o[1]*m[3] + o[2]*m[6]
	self.elements[1] = o[0]*m[1] + o[1]*m[4] + o[2]*m[7]
	self.elements[2] = o[0]*m[2] + o[1]*m[5] + o[2]*m[8]
	self.elements[3] = o[3]*m[0] + o[4]*m[3] + o[5]*m[6] // [0], [1], [2]
	self.elements[4] = o[3]*m[1] + o[4]*m[4] + o[5]*m[7] // [3], [4], [5]
	self.elements[5] = o[3]*m[2] + o[4]*m[5] + o[5]*m[8] // [6], [7], [8]
	self.elements[6] = o[6]*m[0] + o[7]*m[3] + o[8]*m[6]
	self.elements[7] = o[6]*m[1] + o[7]*m[4] + o[8]*m[7]
	self.elements[8] = o[6]*m[2] + o[7]*m[5] + o[8]*m[8]
	return self
}

func (self *Matrix3) SetMultiplyMatrices(matrices ...*Matrix3) *Matrix3 {
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

func (self *Matrix3) Transpose() *Matrix3 {
	o := &self.elements // reference
	return &Matrix3{elements: [9]float32{o[0], o[3], o[6], o[1], o[4], o[7], o[2], o[5], o[8]}}
}

func (self *Matrix3) MultiplyRight(matrix *Matrix3) *Matrix3 {
	o := &self.elements   // reference
	m := &matrix.elements // reference
	return &Matrix3{elements: [9]float32{
		o[0]*m[0] + o[1]*m[3] + o[2]*m[6],
		o[0]*m[1] + o[1]*m[4] + o[2]*m[7],
		o[0]*m[2] + o[1]*m[5] + o[2]*m[8],
		o[3]*m[0] + o[4]*m[3] + o[5]*m[6],
		o[3]*m[1] + o[4]*m[4] + o[5]*m[7],
		o[3]*m[2] + o[4]*m[5] + o[5]*m[8],
		o[6]*m[0] + o[7]*m[3] + o[8]*m[6],
		o[6]*m[1] + o[7]*m[4] + o[8]*m[7],
		o[6]*m[2] + o[7]*m[5] + o[8]*m[8]}}
}
