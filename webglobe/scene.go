package webglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

type Scene struct {
	globe   *SceneGlobe            // the Globe
	objects []*webgl3d.SceneObject // the other SceneObjects
}

func NewScene(wctx *common.WebGLContext) *Scene {
	scene := Scene{}
	scene.globe = NewSceneGlobe(wctx)
	scene.objects = make([]*webgl3d.SceneObject, 0)
	return &scene
}

// ----------------------------------------------------------------------------
// Handling SceneObject
// ----------------------------------------------------------------------------

func (self *Scene) Add(sobj *webgl3d.SceneObject) *Scene {
	if sobj != nil {
		self.objects = append(self.objects, sobj)
	}
	return self
}

func (self *Scene) Get(indices ...int) *webgl3d.SceneObject {
	// Find a SceneObject using the list of indices
	// (multiple indices refers to children[i] of SceneObject)
	scene_object_list := self.objects
	for i := 0; i < len(indices); i++ {
		index := indices[i]
		if index < 0 || index >= len(scene_object_list) {
			return nil
		} else if i == len(indices)-1 {
			scene_object := scene_object_list[index]
			return scene_object
		} else {
			scene_object := scene_object_list[index]
			scene_object_list = scene_object.GetChildren()
		}
	}
	return nil
}
