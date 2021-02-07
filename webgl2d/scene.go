package webgl2d

type Scene struct {
	Objects []*SceneObject
}

func NewScene() *Scene {
	var scene Scene
	scene.Objects = make([]*SceneObject, 0)
	return &scene
}

func (self *Scene) Add(sobj *SceneObject) *Scene {
	self.Objects = append(self.Objects, sobj)
	return self
}
