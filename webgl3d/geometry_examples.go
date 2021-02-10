package webgl3d

func (self *Geometry) LoadCube(xsize float32, ysize float32, zsize float32) *Geometry {
	self.Clear(true, true, true)
	return self
}

func (self *Geometry) LoadSphere(radius float32, wsegs int, hsegs int) *Geometry {
	self.Clear(true, true, true)
	return self
}

func (self *Geometry) LoadCylinder(radius float32, segs int, height float32) *Geometry {
	self.Clear(true, true, true)
	return self
}

func (self *Geometry) LoadPyramid(radius float32, segs int, height float32) *Geometry {
	self.Clear(true, true, true)
	return self
}
