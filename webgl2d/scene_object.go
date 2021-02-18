package webgl2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneObject struct {
	geometry        *Geometry         // geometry
	material        *Material         // material
	shader          *common.Shader    // shader and its bindings
	modelmatrix     geom2d.Matrix3    // model transformation matrix of this SceneObject
	poses           *SceneObjectPoses // poses for multiple instances of this (geometry+material) object
	children        []*SceneObject    // children of this SceneObject (to be rendered recursively)
	parent_material *Material         // shader of the parent SceneObject
	parent_shader   *common.Shader    // shader of the parent SceneObject
	bbox            [2][2]float32     // bounding box
}

func NewSceneObject(geometry *Geometry, material *Material, shader *common.Shader) *SceneObject {
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sobj := SceneObject{geometry: geometry, material: material, shader: shader}
	sobj.modelmatrix.SetIdentity()
	sobj.poses = nil
	sobj.children = nil
	sobj.parent_material = nil
	sobj.parent_shader = nil
	sobj.bbox = geom2d.BBoxInit()
	return &sobj
}

func (self *SceneObject) SetInstancePoses(poses *SceneObjectPoses) *SceneObject {
	self.poses = poses
	return self
}

func (self *SceneObject) AddChild(child *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	if self.material != nil {
		child.parent_material = self.material
	} else {
		child.parent_material = self.parent_material
	}
	if self.shader != nil {
		child.parent_shader = self.shader
	} else {
		child.parent_shader = self.parent_shader
	}
	self.children = append(self.children, child)
	return self
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.geometry.ShowInfo()
	if self.poses != nil {
		fmt.Printf("SceneObject ")
		self.poses.ShowInfo()
	}
	fmt.Printf("SceneObject ")
	self.material.ShowInfo()
	fmt.Printf("SceneObject ")
	self.shader.ShowInfo()
	fmt.Printf("SceneObject children : %d\n", len(self.children))
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txy [2]float32, angle_in_degree float32, sxy [2]float32) *SceneObject {
	translation := geom2d.NewMatrix3().Set(
		1.0, 0.0, txy[0],
		0.0, 1.0, txy[1],
		0.0, 0.0, 1.0)
	radian := float64(angle_in_degree) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	rotation := geom2d.NewMatrix3().Set(
		cos, -sin, 0.0,
		+sin, cos, 0.0,
		0.0, 0.0, 1.0)
	scaling := geom2d.NewMatrix3().Set(
		sxy[0], 0.0, 0.0,
		0.0, sxy[1], 0.0,
		0.0, 0.0, 1.0)
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	radian := float64(angle_in_degree) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	rotation := geom2d.NewMatrix3().Set(
		cos, -sin, 0.0,
		+sin, cos, 0.0,
		0.0, 0.0, 1.0)
	self.modelmatrix = *rotation.MultiplyRight(&self.modelmatrix)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	translation := geom2d.NewMatrix3().Set(
		1.0, 0.0, tx,
		0.0, 1.0, ty,
		0.0, 0.0, 1.0)
	self.modelmatrix = *translation.MultiplyRight(&self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32) *SceneObject {
	scaling := geom2d.NewMatrix3().Set(
		sx, 0.0, 0.0,
		0.0, sy, 0.0,
		0.0, 0.0, 1.0)
	self.modelmatrix = *scaling.MultiplyRight(&self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func (self *SceneObject) GetBoundingBox(m *geom2d.Matrix3, renew bool) [2][2]float32 {
	if !geom2d.BBoxIsSet(self.bbox) || renew {
		bbox := geom2d.BBoxInit()
		// apply the transformation matrx
		var mm *geom2d.Matrix3 = nil
		if m != nil {
			mm = m.MultiplyRight(&self.modelmatrix)
		} else {
			mm = self.modelmatrix.Copy()
		}
		// add all the vertices of the geometry
		if self.poses == nil {
			for _, v := range self.geometry.verts {
				xy := mm.MultiplyVector2(v)
				geom2d.BBoxAddPoint(&bbox, xy)
			}
		} else {
			for i := 0; i < self.poses.count; i++ {
				idx := i * self.poses.size
				txy := self.poses.data_buffer[idx : idx+2]
				for _, v := range self.geometry.verts {
					xy := [2]float32{v[0] + txy[0], v[1] + txy[1]}
					xy = mm.MultiplyVector2(xy)
					geom2d.BBoxAddPoint(&bbox, xy)
				}
			}
		}
		for _, sobj := range self.children {
			bbox = geom2d.BBoxMerge(bbox, sobj.GetBoundingBox(mm, renew))
		}
		self.bbox = bbox
	}
	return self.bbox
}
