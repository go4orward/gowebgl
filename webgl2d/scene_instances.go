package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneInstances struct {
	geometry *Geometry      //
	poses    [][9]float32   //
	colors   [][3]float32   //
	shader   *common.Shader // shader and its bindings

	translation [2]float32 // translation in world space
	rotation    float32    // rotation angle in degree
	scaling     [2]float32 // scaling of the geometry
	modelmatrix *geom2d.Matrix3
}

func NewSceneInstances(geometry *Geometry, poses [][]float32, shader *common.Shader) *SceneInstances {
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sins := SceneInstances{geometry: geometry, poses: poses, colors: colors, shader: shader}
	sins.modelmatrix = geom2d.NewMatrix3() // identity
	sins.translation = [2]float32{0, 0}
	sins.rotation = 0.0
	sins.scaling = [2]float32{1, 1}
	return &sins
}
