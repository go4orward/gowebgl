package webgl2d

import (
	"fmt"

	"github.com/go4orward/gowebgl/geom2d"
	"github.com/go4orward/gowebgl/wcommon"
)

type SceneObject struct {
	Geometry    *Geometry                 // geometry interface
	Material    *wcommon.Material         // material
	VShader     *wcommon.Shader           // vert shader and its bindings
	EShader     *wcommon.Shader           // edge shader and its bindings
	FShader     *wcommon.Shader           // face shader and its bindings
	modelmatrix geom2d.Matrix3            // model transformation matrix of this SceneObject
	UseDepth    bool                      // depth test flag (default is true)
	UseBlend    bool                      // blending flag with alpha (default is false)
	poses       *wcommon.SceneObjectPoses // OPTIONAL, poses for multiple instances of this (geometry+material) object
	children    []*SceneObject            // OPTIONAL, children of this SceneObject (to be rendered recursively)
	bbox        [2][2]float32             // bounding box
}

func NewSceneObject(geometry *Geometry, material *wcommon.Material,
	vshader *wcommon.Shader, eshader *wcommon.Shader, fshader *wcommon.Shader) *SceneObject {
	// 'geometry' : geometric shape (vertices, edges, faces) to be rendered
	// 'material' : color, texture, or other material properties	: OPTIONAL (can be 'nil')
	// 'vshader' : shader for VERTICES (POINTS) 					: OPTIONAL (can be 'nil')
	// 'eshader' : shader for EDGES (LINES) 						: OPTIONAL (can be 'nil')
	// 'fshader' : shader for FACES (TRIANGLES) 					: OPTIONAL (can be 'nil')
	// Note that geometry & material & shader can be shared among different SceneObjects.
	if geometry == nil {
		return nil
	}
	sobj := SceneObject{Geometry: geometry, Material: material, VShader: vshader, EShader: eshader, FShader: fshader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = false // new drawings will overwrite old ones by default
	sobj.UseBlend = false // alpha blending is turned off by default
	sobj.poses = nil      // OPTIONAL, only if multiple instances of the geometry are rendered
	sobj.children = nil   // OPTIONAL, only if current SceneObject has any child SceneObjects
	sobj.bbox = geom2d.BBoxInit()
	return &sobj
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.Geometry.ShowInfo()
	if self.poses != nil {
		fmt.Printf("  ")
		self.poses.ShowInfo()
	}
	if self.Material != nil {
		fmt.Printf("  ")
		self.Material.ShowInfo()
	}
	if self.VShader != nil {
		fmt.Printf("  VERT ")
		self.VShader.ShowInfo()
	}
	if self.EShader != nil {
		fmt.Printf("  EDGE ")
		self.EShader.ShowInfo()
	}
	if self.FShader != nil {
		fmt.Printf("  FACE ")
		self.FShader.ShowInfo()
	}
	fmt.Printf("  Flags    : UseDepth=%t  UseBlend=%t\n", self.UseDepth, self.UseBlend)
	fmt.Printf("  Children : %d\n", len(self.children))
}

// ----------------------------------------------------------------------------
// Basic Access
// ----------------------------------------------------------------------------

func (self *SceneObject) AddChild(child *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	self.children = append(self.children, child)
	return self
}

// ----------------------------------------------------------------------------
// Multiple Instance Poses
// ----------------------------------------------------------------------------

func (self *SceneObject) SetPoses(poses *wcommon.SceneObjectPoses) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses = poses
	return self
}

func (self *SceneObject) SetupPoses(size int, count int, data []float32) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses = wcommon.NewSceneObjectPoses(size, count, data)
	return self
}

func (self *SceneObject) SetPoseValues(index int, offset int, values ...float32) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses.SetPose(index, offset, values...)
	return self
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txy [2]float32, angle_in_degree float32, sxy [2]float32) *SceneObject {
	translation := geom2d.NewMatrix3().SetTranslation(txy[0], txy[1])
	rotation := geom2d.NewMatrix3().SetRotation(angle_in_degree)
	scaling := geom2d.NewMatrix3().SetScaling(sxy[0], sxy[1])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	rotation := geom2d.NewMatrix3().SetRotation(angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	translation := geom2d.NewMatrix3().SetTranslation(tx, ty)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32) *SceneObject {
	scaling := geom2d.NewMatrix3().SetScaling(sx, sy)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func (self *SceneObject) GetBoundingBox(m *geom2d.Matrix3, renew bool) [2][2]float32 {
	if !geom2d.BBoxIsSet(self.bbox) || renew {
		bbox := geom2d.BBoxInit()
		// apply the transformation matrx
		var mm *geom2d.Matrix3 = nil
		if m != nil {
			mm = m.MultiplyToTheRight(&self.modelmatrix)
		} else {
			mm = self.modelmatrix.Copy() // new matrix
		}
		// add all the vertices of the geometry
		if self.poses == nil {
			for _, v := range self.Geometry.verts {
				xy := mm.MultiplyVector2(v)
				geom2d.BBoxAddPoint(&bbox, xy)
			}
		} else {
			for i := 0; i < self.poses.Count; i++ {
				idx := i * self.poses.Size
				txy := self.poses.DataBuffer[idx : idx+2]
				for _, v := range self.Geometry.verts {
					xy := [2]float32{v[0] + txy[0], v[1] + txy[1]}
					xy = mm.MultiplyVector2(xy)
					geom2d.BBoxAddPoint(&bbox, xy)
				}
			}
		}
		for _, sobj := range self.children {
			bbox = geom2d.BBoxMerge(bbox, sobj.GetBoundingBox(mm, renew))
		}
		self.bbox = bbox
	}
	return self.bbox
}
