package webgl2d

import (
	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneObject struct {
	geometry *Geometry
	material *Material
	shader   *Shader // shader program and its bindings

	children    []*SceneObject
	modelMatrix *geom2d.Matrix3
}

func NewSceneObject(geometry *Geometry, material *Material, shader *Shader) *SceneObject {
	sobj := SceneObject{geometry: geometry, material: material, shader: shader}
	sobj.children = nil
	sobj.modelMatrix = nil
	return &sobj
}

func (self *SceneObject) AddChild(sobj *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	self.children = append(self.children, sobj)
	return self
}
