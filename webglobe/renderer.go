package webglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

type Renderer struct {
	wctx     *common.WebGLContext
	render3d *webgl3d.Renderer
	axes     *webgl3d.SceneObject
}

func NewRenderer(wctx *common.WebGLContext) *Renderer {
	renderer := Renderer{wctx: wctx, render3d: webgl3d.NewRenderer(wctx), axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *Renderer) Clear(wcamera *WorldCamera, color string) {
	self.render3d.Clear(wcamera.gcam, color)
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *Renderer) RenderAxes(wcamera *WorldCamera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	self.render3d.RenderAxes(wcamera.gcam, length)
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(wcamera *WorldCamera, scene *Scene) {
	// Render the globe
	new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(scene.Globe.globe_obj.GetModelMatrix())
	self.render3d.RenderSceneObject(scene.Globe.globe_obj, wcamera.gcam.GetProjMatrix(), new_viewmodel)

	// Render all the SceneObjects in the Scene
	for _, sobj := range scene.objects {
		new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(sobj.GetModelMatrix())
		self.render3d.RenderSceneObject(sobj, wcamera.gcam.GetProjMatrix(), new_viewmodel)
	}
}
