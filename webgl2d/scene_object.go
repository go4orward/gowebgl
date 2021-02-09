package webgl2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneObject struct {
	geometry        *Geometry
	material        *Material // material
	shader          *Shader   // shader and its bindings
	parent_material *Material // shader of the parent SceneObject
	parent_shader   *Shader   // shader of the parent SceneObject

	translation []float32 // translation in world space
	rotation    float32   // rotation angle in degree
	scaling     float32   // scaling of the geometry
	modelmatrix *geom2d.Matrix3
	children    []*SceneObject
}

func NewSceneObject(geometry *Geometry, material *Material, shader *Shader) *SceneObject {
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sobj := SceneObject{geometry: geometry, material: material, shader: shader}
	sobj.modelmatrix = geom2d.NewMatrix3() // identity
	sobj.translation = []float32{0, 0}
	sobj.rotation = 0.0
	sobj.scaling = 1.0
	sobj.children = nil
	sobj.parent_material = nil
	sobj.parent_shader = nil
	return &sobj
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

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	self.rotation = angle_in_degree
	return self.update_model_matrix()
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	self.translation[0] = tx
	self.translation[1] = ty
	return self.update_model_matrix()
}

func (self *SceneObject) Scale(scale float32) *SceneObject {
	self.scaling = scale
	return self.update_model_matrix()
}

func (self *SceneObject) Transform(txy []float32, angle_in_degree float32, scale float32) *SceneObject {
	if scale <= 0.0 || scale >= 1000.0 {
		fmt.Printf("SceneObject.Transform() failed : invalid scale = %.1f\n", scale)
		return self
	}
	geom2d.Assign(self.translation, txy)
	self.rotation = angle_in_degree
	self.scaling = scale
	return self.update_model_matrix()
}

func (self *SceneObject) update_model_matrix() *SceneObject {
	radian := float64(self.rotation) * (math.Pi / 180.0)
	cos, sin := float32(math.Cos(radian)), float32(math.Sin(radian))
	translation := geom2d.NewMatrix3().Set(
		1.0, 0.0, self.translation[0],
		0.0, 1.0, self.translation[1],
		0.0, 0.0, 1.0)
	rotation := geom2d.NewMatrix3().Set(
		cos, -sin, 0.0,
		sin, cos, 0.0,
		0.0, 0.0, 1.0)
	scaling := geom2d.NewMatrix3().Set(
		self.scaling, 0.0, 0.0,
		0.0, self.scaling, 0.0,
		0.0, 0.0, 1.0)
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}
