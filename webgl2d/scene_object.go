package webgl2d

import (
	"fmt"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneObject struct {
	Geometry    *Geometry         // geometry
	Material    *Material         // material
	Shader      *common.Shader    // shader and its bindings
	modelmatrix geom2d.Matrix3    // model transformation matrix of this SceneObject
	UseDepth    bool              // depth test flag (default is true)
	UseBlend    bool              // blending flag with alpha (default is false)
	poses       *SceneObjectPoses // poses for multiple instances of this (geometry+material) object
	children    []*SceneObject    // children of this SceneObject (to be rendered recursively)
	bbox        [2][2]float32     // bounding box
}

func NewSceneObject(geometry *Geometry, material *Material, shader *common.Shader) *SceneObject {
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sobj := SceneObject{Geometry: geometry, Material: material, Shader: shader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = true
	sobj.UseBlend = false
	sobj.poses = nil
	sobj.children = nil
	sobj.bbox = geom2d.BBoxInit()
	return &sobj
}

// ----------------------------------------------------------------------------
// Basic Access
// ----------------------------------------------------------------------------

func (self *SceneObject) SetInstancePoses(poses *SceneObjectPoses) *SceneObject {
	self.poses = poses
	return self
}

func (self *SceneObject) AddChild(child *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	self.children = append(self.children, child)
	return self
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.Geometry.ShowInfo()
	if self.poses != nil {
		fmt.Printf("SceneObject ")
		self.poses.ShowInfo()
	}
	fmt.Printf("SceneObject ")
	self.Material.ShowInfo()
	fmt.Printf("SceneObject ")
	self.Shader.ShowInfo()
	fmt.Printf("SceneObject children : %d\n", len(self.children))
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txy [2]float32, angle_in_degree float32, sxy [2]float32) *SceneObject {
	translation := geom2d.NewMatrix3().SetTranslation(txy[0], txy[1])
	rotation := geom2d.NewMatrix3().SetRotation(angle_in_degree)
	scaling := geom2d.NewMatrix3().SetScaling(sxy[0], sxy[1])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	rotation := geom2d.NewMatrix3().SetRotation(angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	translation := geom2d.NewMatrix3().SetTranslation(tx, ty)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32) *SceneObject {
	scaling := geom2d.NewMatrix3().SetScaling(sx, sy)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
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
			mm = m.MultiplyToTheRight(&self.modelmatrix)
		} else {
			mm = self.modelmatrix.Copy() // new matrix
		}
		// add all the vertices of the geometry
		if self.poses == nil {
			for _, v := range self.Geometry.verts {
				xy := mm.MultiplyVector2(v)
				geom2d.BBoxAddPoint(&bbox, xy)
			}
		} else {
			for i := 0; i < self.poses.count; i++ {
				idx := i * self.poses.size
				txy := self.poses.data_buffer[idx : idx+2]
				for _, v := range self.Geometry.verts {
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
