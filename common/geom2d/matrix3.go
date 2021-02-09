package geom2d

type Matrix3 struct {
	elements []float32
}

func NewMatrix3() *Matrix3 {
	mtx := Matrix3{elements: []float32{1, 0, 0, 0, 1, 0, 0, 0, 1}} // identity matrix
	return &mtx
}

func (self *Matrix3) GetElements() []float32 {
	return self.elements
}

// ----------------------------------------------------------------------------
// Setting element values
// ----------------------------------------------------------------------------

func (self *Matrix3) Set(v00 float32, v01 float32, v02 float32, v10 float32, v11 float32, v12 float32, v20 float32, v21 float32, v22 float32) *Matrix3 {
	self.elements[0] = v00
	self.elements[1] = v01
	self.elements[2] = v02
	self.elements[3] = v10
	self.elements[4] = v11
	self.elements[5] = v12
	self.elements[6] = v20
	self.elements[7] = v21
	self.elements[8] = v22
	return self
}

func (self *Matrix3) SetTranspose() *Matrix3 {
	self.elements[1], self.elements[3] = self.elements[3], self.elements[1] // [0], [1], [2]
	self.elements[2], self.elements[6] = self.elements[6], self.elements[2] // [3], [4], [5]
	self.elements[5], self.elements[7] = self.elements[7], self.elements[5] // [6], [7], [8]
	return self
}

func (self *Matrix3) SetCopy(m *Matrix3) *Matrix3 {
	copy(self.elements, m.elements)
	return self
}

func (self *Matrix3) SetMultiplyRight(matrix *Matrix3) *Matrix3 {
	o := []float32{0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(o, self.elements)
	m := matrix.elements
	e := self.elements
	e[0] = o[0]*m[0] + o[1]*m[3] + o[2]*m[6]
	e[1] = o[0]*m[1] + o[1]*m[4] + o[2]*m[7]
	e[2] = o[0]*m[2] + o[1]*m[5] + o[2]*m[8]
	e[3] = o[3]*m[0] + o[4]*m[3] + o[5]*m[6]
	e[4] = o[3]*m[1] + o[4]*m[4] + o[5]*m[7]
	e[5] = o[3]*m[2] + o[4]*m[5] + o[5]*m[8]
	e[6] = o[6]*m[0] + o[7]*m[3] + o[8]*m[6]
	e[7] = o[6]*m[1] + o[7]*m[4] + o[8]*m[7]
	e[8] = o[6]*m[2] + o[7]*m[5] + o[8]*m[8]
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
	o := self.elements
	elements := []float32{
		o[0], o[3], o[6],
		o[1], o[4], o[7],
		o[2], o[5], o[8]}
	return &Matrix3{elements: elements}
}

func (self *Matrix3) MultiplyRight(matrix *Matrix3) *Matrix3 {
	o := self.elements
	m := matrix.elements
	elements := []float32{
		o[0]*m[0] + o[1]*m[3] + o[2]*m[6],
		o[0]*m[1] + o[1]*m[4] + o[2]*m[7],
		o[0]*m[2] + o[1]*m[5] + o[2]*m[8],
		o[3]*m[0] + o[4]*m[3] + o[5]*m[6],
		o[3]*m[1] + o[4]*m[4] + o[5]*m[7],
		o[3]*m[2] + o[4]*m[5] + o[5]*m[8],
		o[6]*m[0] + o[7]*m[3] + o[8]*m[6],
		o[6]*m[1] + o[7]*m[4] + o[8]*m[7],
		o[6]*m[2] + o[7]*m[5] + o[8]*m[8]}
	return &Matrix3{elements: elements}
}
