package webglobe

import (
	"image"
	"math"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/webgl3d"
)

type Globe struct {
	globe_obj *webgl3d.SceneObject
	// TODO: Add glow ring around the globe profile
}

func NewGlobe(wctx *common.WebGLContext) *Globe {
	globe := Globe{}
	globe.globe_obj = NewSceneObject_Globe(wctx) // texture & vertex normals
	// log.Println("Please wait while world image is loaded.") // printed in the browser console
	return &globe
}

func (self *Globe) Rotate(angle_in_degree float32) {
	self.globe_obj.Rotate([3]float32{0, 0, 1}, angle_in_degree)
}

// ----------------------------------------------------------------------------
// Globe
// ----------------------------------------------------------------------------

func NewSceneObject_GlobeWithoutLight(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// This Globe model has texture only (without normal vectors and directional lighting).
	geometry := NewGeometry_Globe(1.0, 64, 32)                 // create geometry (a cube of size 1.0)
	geometry.BuildDataBuffers(true, false, true)               // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.png") // create material with a texture of world image
	shader := webgl3d.NewShader_Texture(wctx)                  // create a shader, and set its bindings
	return webgl3d.NewSceneObject(geometry, material, shader)  // set up the scene object
}

func NewSceneObject_Globe(wctx *common.WebGLContext) *webgl3d.SceneObject {
	// This Globe model has texture AND normal vectors (with directional lighting).
	geometry := NewGeometry_Globe(1.0, 64, 32) // create geometry (a cube of size 1.0)
	geometry.BuildNormalsForVertex()           // calculate normal vectors for each vertex
	for i := 0; i <= 64; i++ {                 // Note that there are multiple north/south pole vertices
		geometry.ChangeNormal(i*33+0, [3]float32{0, 0, -1})  // adjust normal vector at south pole
		geometry.ChangeNormal(i*33+32, [3]float32{0, 0, +1}) // adjust normal vector at north pole
	}
	geometry.BuildDataBuffers(true, false, true)               // build data buffers for vertices and faces
	material := webgl3d.NewMaterial(wctx, "/assets/world.png") // create material (yellow color)
	shader := webgl3d.NewShader_BasicTexture(wctx)             // create a shader, and set its bindings
	return webgl3d.NewSceneObject(geometry, material, shader)  // set up the scene object
}

func NewGeometry_Globe(radius float32, wsegs int, hsegs int) *webgl3d.Geometry {
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
		}
	}
	for i := 0; i < wnum-1; i++ { // quadratic faces on the side (triangles for south/north poles)
		for j := 0; j < hnum-1; j++ {
			start := uint32((i+0)*hnum + j)
			wnext := uint32((i+1)*hnum + j)
			if spole := (j == 0); spole {
				geometry.AddFace([]uint32{start, wnext + 1, start + 1})
			} else if npole := (j == hsegs-1); npole {
				geometry.AddFace([]uint32{start, wnext + 0, start + 1})
			} else {
				geometry.AddFace([]uint32{start, wnext, wnext + 1, start + 1})
			}
			// fmt.Printf("adding Face [ %d %d %d %d ]\n", start, wnext, wnext+1, start+1)
		}
	}
	return geometry
}

// ----------------------------------------------------------------------------
// GlowRing around the Globe
// ----------------------------------------------------------------------------

func NewGeometry_GlowRing(in_radius float32, out_radius float32, nsegs int, z float32) *webgl3d.Geometry {
	geometry := webgl3d.NewGeometry()
	rad := math.Pi * 2 / float64(nsegs)
	for i := 0; i < nsegs; i++ {
		cos, sin := float32(math.Cos(rad*float64(i))), float32(math.Sin(rad*float64(i)))
		geometry.AddVertex([3]float32{in_radius * cos, in_radius * sin, z})
		geometry.AddVertex([3]float32{out_radius * cos, out_radius * sin, z})
		geometry.AddTextureUV([]float32{0.0, 1.0})
		geometry.AddTextureUV([]float32{1.0, 1.0})
		ii, jj := uint32(i), uint32((i+1)%nsegs)
		geometry.AddFace([]uint32{2*ii + 0, 2*jj + 0, 2*jj + 1, 2*ii + 1})
	}
	return geometry
}

func NewTextureImage_Glow() *image.RGBA {
	size, color := 32, [4]float32{0.8, 0.8, 0.8, 1.0}
	pbuffer := make([]uint8, size*4)
	for i, pos := 0, 0; i < size; i++ {
		s := 1.0 - (float32(i) / float32(size-1))
		pbuffer[pos+0] = uint8(color[0] * 255.0)
		pbuffer[pos+1] = uint8(color[1] * 255.0)
		pbuffer[pos+2] = uint8(color[2] * 255.0)
		pbuffer[pos+4] = uint8(color[3] * 255.0 * s * s)
		pos += 4
	}
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{32, 3}})
	return img
}

// public static loadGlowTexture(unit: number, size: number, color: string): Material {
// 	// Glow texture for highlighted features
// 	//   v == 0.0 : ascending  by u [0.0 ~ 0.5 ~ 1.0]
// 	//   v == 0.5 : both_side       [0.0 ~ 1.0 ~ 0.0]
// 	//   v == 1.0 : descending by u [1.0 ~ 0.5 ~ 0.0]
// 	let rgba = Color.getRGBA(color);
// 	let data = new Uint8Array(size * 4 * 3);
// 	let idx0 = size * 4 * 0, idx1 = size * 4 * 1, idx2 = size * 4 * 2;
// 	for (let u = 0; u < size; u++) {
// 		// ascending glow for the first  row     (v == 0.0)  [ 0.0 ~ 0.5 ~ 1.0 ]
// 		let intensity = (u / (size - 1));
// 		intensity = intensity * intensity;
// 		data[idx0+0] = rgba[0]*0xff;  data[idx0+1] = rgba[1]*0xff;  data[idx0+2] = rgba[2]*0xff;
// 		data[idx0+3] = intensity * rgba[3] * 0xff;
// 		idx0 += 4;
// 		// glow on both side for the second row  (v == 0.5)  [ 0.0 ~ 1.0 ~ 0.0 ]
// 		intensity = 1.0 - Math.abs((u / (size - 1)) * 2 - 1);
// 		intensity = intensity * intensity;
// 		data[idx1+0] = rgba[0]*0xff;  data[idx1+1] = rgba[1]*0xff;  data[idx1+2] = rgba[2]*0xff;
// 		data[idx1+3] = intensity * rgba[3] * 0xff;
// 		idx1 += 4;
// 		// descending glow for the third row     (v == 1.0)  [ 1.0 ~ 0.5 ~ 0.0 ]
// 		intensity = 1.0 - (u / (size - 1));
// 		intensity = intensity * intensity;
// 		data[idx2+0] = rgba[0]*0xff;  data[idx2+1] = rgba[1]*0xff;  data[idx2+2] = rgba[2]*0xff;
// 		data[idx2+3] = intensity * rgba[3] * 0xff;
// 		idx2 += 4;
// 	}
// 	let material = new Material().setTexture(unit, data, size, 3, TextureFormat.RGBA, "glow");
// 	material.setTextureFilters(TextureFilter.NEAREST, TextureFilter.NEAREST, TextureFilter.CLAMP_TO_EDGE, TextureFilter.CLAMP_TO_EDGE);
// 	return material;
// }
