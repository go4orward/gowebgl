package webgl2d

import "github.com/go4orward/gowebgl/common/geom2d"

type Scene struct {
	objects []*SceneObject // list of SceneObjects in the Scene
	bbox    [2][2]float32  // bounding box of all the SceneObjects
}

func NewScene() *Scene {
	var scene Scene
	scene.objects = make([]*SceneObject, 0)
	scene.bbox = geom2d.BBoxInit()
	return &scene
}

// ----------------------------------------------------------------------------
// Handling SceneObject
// ----------------------------------------------------------------------------

func (self *Scene) Add(sobj *SceneObject) *Scene {
	if sobj != nil {
		self.objects = append(self.objects, sobj)
	}
	return self
}

func (self *Scene) Get(indices ...int) *SceneObject {
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
			scene_object_list = scene_object.children
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func (self *Scene) GetBoundingBox(renew bool) [2][2]float32 {
	if !geom2d.BBoxIsSet(self.bbox) || renew {
		bbox := geom2d.BBoxInit()
		for _, sobj := range self.objects {
			bbox = geom2d.BBoxMerge(bbox, sobj.GetBoundingBox(nil, renew))
		}
		self.bbox = bbox
	}
	return self.bbox
}

func (self *Scene) GetBBoxSizeCenter(renew bool) ([2][2]float32, [2]float32, [2]float32) {
	bbox := self.GetBoundingBox(renew)
	return bbox, geom2d.BBoxSize(bbox), geom2d.BBoxCenter(bbox)
}
