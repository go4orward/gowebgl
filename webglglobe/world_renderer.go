package webglglobe

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom3d"
	"github.com/go4orward/gowebgl/webgl3d"
)

type WorldRenderer struct {
	wctx     *common.WebGLContext
	renderer *webgl3d.Renderer
	axes     *webgl3d.SceneObject
}

func NewWorldRenderer(wctx *common.WebGLContext) *WorldRenderer {
	renderer := WorldRenderer{wctx: wctx, renderer: webgl3d.NewRenderer(wctx), axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *WorldRenderer) Clear(wcamera *WorldCamera, color string) {
	self.renderer.Clear(wcamera.gcam, color)
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderAxes(wcamera *WorldCamera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	self.renderer.RenderAxes(wcamera.gcam, length)
}

// ----------------------------------------------------------------------------
// Rendering the World
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderWorld(wcamera *WorldCamera, globe *Globe) {
	if globe.IsReadyToRender() {
		// Render the Globe
		new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(&globe.modelmatrix)
		self.renderer.RenderSceneObject(globe.gsphere, wcamera.gcam.GetProjMatrix(), new_viewmodel)
		// Render the GlowRing (in CAMERA space)
		distance := geom3d.Length(wcamera.gcam.GetCenter())
		translation := geom3d.NewMatrix4().SetTranslation(0, 0, -distance)
		self.renderer.RenderSceneObject(globe.glowring, wcamera.gcam.GetProjMatrix(), translation)
	}
}
