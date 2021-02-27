package webgl3d

import (
	"fmt"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom3d"
)

type SceneObject struct {
	Geometry    *Geometry
	Material    *Material         // material
	Shader      *common.Shader    // shader and its bindings
	modelmatrix geom3d.Matrix4    //
	UseDepth    bool              // depth test flag (default is true)
	UseBlend    bool              // blending flag with alpha (default is false)
	poses       *SceneObjectPoses // poses for multiple instances of this (geometry+material) object
	children    []*SceneObject    //
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

func (self *SceneObject) GetModelMatrix() *geom3d.Matrix4 {
	return &self.modelmatrix
}

func (self *SceneObject) GetChildren() []*SceneObject {
	return self.children
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.Geometry.ShowInfo()
	fmt.Printf("SceneObject ")
	self.Material.ShowInfo()
	fmt.Printf("SceneObject ")
	self.Shader.ShowInfo()
	fmt.Printf("SceneObject children : %d\n", len(self.children))
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txyz [3]float32, axis [3]float32, angle_in_degree float32, sxyz [3]float32) *SceneObject {
	translation := geom3d.NewMatrix4().SetTranslation(txyz[0], txyz[1], txyz[2])
	rotation := geom3d.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	scaling := geom3d.NewMatrix4().SetScaling(sxyz[0], sxyz[1], sxyz[2])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32, tz float32) *SceneObject {
	translation := geom3d.NewMatrix4().SetTranslation(tx, ty, tz)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Rotate(axis [3]float32, angle_in_degree float32) *SceneObject {
	rotation := geom3d.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32, sz float32) *SceneObject {
	scaling := geom3d.NewMatrix4().SetScaling(sx, sy, sz)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}
