package webgl3d

import "github.com/go4orward/gowebgl/wcommon"

type Scene struct {
	bkgcolor [3]float32     // background color of the scene
	objects  []*SceneObject // SceneObjects in the scene
	overlays []Overlay      // list of Overlay (interface) layers
}

func NewScene(bkg_color string) *Scene {
	var scene Scene
	scene.SetBkgColor(bkg_color)
	scene.objects = make([]*SceneObject, 0)
	scene.overlays = make([]Overlay, 0)
	return &scene
}

// ----------------------------------------------------------------------------
// Background Color
// ----------------------------------------------------------------------------

func (self *Scene) SetBkgColor(color string) *Scene {
	rgba := wcommon.ParseHexColor(color)
	self.bkgcolor = [3]float32{rgba[0], rgba[1], rgba[2]}
	return self
}

func (self *Scene) GetBkgColor() [3]float32 {
	return self.bkgcolor
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
// Managing OverlayLayers
// ----------------------------------------------------------------------------

func (self *Scene) AddOverlay(overlay ...Overlay) *Scene {
	for i := 0; i < len(overlay); i++ {
		self.overlays = append(self.overlays, overlay[i])
	}
	return self
}
