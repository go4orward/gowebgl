package webglglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

type WorldRenderer struct {
	wctx     *common.WebGLContext
	render3d *webgl3d.Renderer
	axes     *webgl3d.SceneObject
}

func NewWorldRenderer(wctx *common.WebGLContext) *WorldRenderer {
	renderer := WorldRenderer{wctx: wctx, render3d: webgl3d.NewRenderer(wctx), axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *WorldRenderer) Clear(wcamera *WorldCamera, color string) {
	self.render3d.Clear(wcamera.gcam, color)
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderAxes(wcamera *WorldCamera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	self.render3d.RenderAxes(wcamera.gcam, length)
}

// ----------------------------------------------------------------------------
// Rendering the World
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderWorld(wcamera *WorldCamera, globe *Globe) {
	// Render the globe
	new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(globe.gsphere.GetModelMatrix())
	self.render3d.RenderSceneObject(globe.gsphere, wcamera.gcam.GetProjMatrix(), new_viewmodel)
}
