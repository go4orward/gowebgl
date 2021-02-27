package webglobe

import (
	"math"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom3d"
	"github.com/go4orward/gowebgl/webgl3d"
)

type Globe struct {
	gsphere  *webgl3d.SceneObject //
	glowring *webgl3d.SceneObject //
	// TODO: Add glow ring around the globe profile
}

func NewGlobe(wctx *common.WebGLContext) *Globe {
	self := Globe{}
	self.gsphere = NewSceneObject_Globe(wctx)     // texture & vertex normals
	self.glowring = NewSceneObject_GlowRing(wctx) //
	return &self
}

func (self *Globe) Rotate(angle_in_degree float32) {
	self.gsphere.Rotate([3]float32{0, 0, 1}, angle_in_degree)
}

// ----------------------------------------------------------------------------
// Globe
// ----------------------------------------------------------------------------

func NewSceneObject_GlobeWithoutLight(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// Globe model with texture UV coordinates (without normal vectors and directional lighting)
	gsphere := build_globe_geometry(1.0, 64, 32, false)        // create globe geometry with texture UVs only
	gsphere.BuildDataBuffers(true, false, true)                // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.png") // create material with a texture of world image
	shader := webgl3d.NewShader_TextureOnly(wctx)              // use the standard TEXTURE_ONLY shader
	return webgl3d.NewSceneObject(gsphere, material, shader)   // set up the scene object
}

func NewSceneObject_Globe(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// Globe model with texture AND normal vectors (for directional lighting)
	gsphere := build_globe_geometry(1.0, 64, 32, true)         // create globe geometry with vertex normal vectors
	gsphere.BuildDataBuffers(true, false, true)                // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.png") // create material with a texture of world image
	shader := webgl3d.NewShader_NormalTexture(wctx)            // use the standard NORMAL+TEXTURE shader
	return webgl3d.NewSceneObject(gsphere, material, shader)   // set up the scene object
}

func build_globe_geometry(radius float32, wsegs int, hsegs int, use_normals bool) *webgl3d.Geometry {
	// Globe (sphere) geometry with UV coordinates per vertex (to be used with a texture image)
	//   Note that multiple vertices are assigned to north/south poles, as well as 0/360 longitude.
	//   This approach results in more efficient data buffers than a simple sphere,
	//   since we can build the buffers with single point per vertex, without any repetition.
	geometry := webgl3d.NewGeometry()
	wnum, hnum := wsegs+1, hsegs+1
	wstep := math.Pi * 2.0 / float32(wsegs)
	hstep := math.Pi / float32(hsegs)
	for i := 0; i < wnum; i++ {
		lon := wstep*float32(i) - math.Pi // longitude (λ) [-180 ~ 180]
		for j := 0; j < hnum; j++ {
			lat := -math.Pi/2.0 + hstep*float32(j) // latitude (φ)
			xyz := GetXYZFromLL(lon, lat, radius)
			geometry.AddVertex(xyz)
			geometry.AddTextureUV([]float32{float32(i) / float32(wsegs), 1.0 - float32(j)/float32(hsegs)})
			if use_normals {
				geometry.AddNormal(geom3d.Normalize(xyz))
			}
		}
	}
	for i := 0; i < wnum-1; i++ { // faces on the side
		for j := 0; j < hnum-1; j++ {
			start := uint32((i+0)*hnum + j)
			wnext := uint32((i+1)*hnum + j)
			if spole := (j == 0); spole {
				geometry.AddFace([]uint32{start, wnext + 1, start + 1}) // triangular face for south pole
			} else if npole := (j == hsegs-1); npole {
				geometry.AddFace([]uint32{start, wnext + 0, start + 1}) // triangular face for north pole
			} else {
				geometry.AddFace([]uint32{start, wnext, wnext + 1, start + 1}) // quadratic face otherwise
			}
		}
	}
	return geometry
}

// ----------------------------------------------------------------------------
// GlowRing around the Globe
// ----------------------------------------------------------------------------

func NewSceneObject_GlowRing(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// GlowRing around the globe, to make the globe stand out against black background.
	// (Note that GlowRing should be rendered in CAMERA space by Renderer)
	glowring := build_glowring_geometry(1.0, 1.2, 64)             // create geometry (a ring around the globe)
	glowring.BuildDataBuffers(true, false, true)                  // build data buffers for vertices and faces
	material := webgl3d.NewMaterialForGlowEffect(wctx, "#ffffff") // texture material for glow effect
	shader := webgl3d.NewShader_TextureOnly(wctx)                 // use the standard TEXTURE_ONLY shader
	return webgl3d.NewSceneObject(glowring, material, shader)     // set up the scene object
}

func build_glowring_geometry(in_radius float32, out_radius float32, nsegs int) *webgl3d.Geometry {
	geometry := webgl3d.NewGeometry()
	rad := math.Pi * 2 / float64(nsegs)
	for i := 0; i < nsegs; i++ {
		cos, sin := float32(math.Cos(rad*float64(i))), float32(math.Sin(rad*float64(i)))
		geometry.AddVertex([3]float32{in_radius * cos, in_radius * sin, 0})
		geometry.AddVertex([3]float32{out_radius * cos, out_radius * sin, 0})
		geometry.AddTextureUV([]float32{0.0, 1.0})
		geometry.AddTextureUV([]float32{1.0, 1.0})
		ii, jj := uint32(i), uint32((i+1)%nsegs)
		geometry.AddFace([]uint32{2*ii + 0, 2*jj + 0, 2*jj + 1, 2*ii + 1})
	}
	return geometry
}
