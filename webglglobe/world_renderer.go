package webglglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom3d"
	"github.com/go4orward/gowebgl/webgl3d"
)

type WorldRenderer struct {
	wctx     *common.WebGLContext // WebGL context
	renderer *webgl3d.Renderer    // Renderer for rendering 3D SceneObjects
	axes     *webgl3d.SceneObject // XYZ axes for visual reference (only if required)
}

func NewWorldRenderer(wctx *common.WebGLContext) *WorldRenderer {
	renderer := WorldRenderer{wctx: wctx, renderer: webgl3d.NewRenderer(wctx), axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *WorldRenderer) Clear(globe *Globe) {
	context := self.wctx.GetContext()
	constants := self.wctx.GetConstants()
	rgb := globe.GetBkgColor()
	context.Call("clearColor", rgb[0], rgb[1], rgb[2], 1.0) // set clearing color
	context.Call("clear", constants.COLOR_BUFFER_BIT)       // clear the canvas
	context.Call("clear", constants.DEPTH_BUFFER_BIT)       // clear the canvas
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderAxes(wcamera *WorldCamera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	if self.axes == nil {
		self.axes = webgl3d.NewSceneObject_3DAxes(self.wctx, length)
	}
	self.renderer.RenderSceneObject(self.axes, wcamera.gcam.GetProjMatrix(), wcamera.gcam.GetViewMatrix())
}

// ----------------------------------------------------------------------------
// Rendering the World
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderWorld(globe *Globe, wcamera *WorldCamera) {
	if globe.IsReadyToRender() {
		// Render the Globe
		// new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(&globe.modelmatrix)
		// self.renderer.RenderSceneObject(globe.GSphere, wcamera.gcam.GetProjMatrix(), new_viewmodel)
		// Render the GlowRing (in CAMERA space)
		distance := geom3d.Length(wcamera.gcam.GetCenter())
		translation := geom3d.NewMatrix4().SetTranslation(0, 0, -distance)
		self.renderer.RenderSceneObject(globe.GlowRing, wcamera.gcam.GetProjMatrix(), translation)
	}
}
