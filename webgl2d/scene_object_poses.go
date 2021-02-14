package webgl2d

import (
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
	"github.com/go4orward/gowebgl/common/geom2d"
)

type SceneObjectPoses struct {
	poses              []geom2d.Matrix3 //
	poses_data_buffer  []float32        //
	poses_webgl_buffer js.Value         //
	poses_info         [3]int           //
}

func NewSceneObjectPoses() *SceneObjectPoses {
	poses := SceneObjectPoses{}
	return &poses
}

// ------------------------------------------------------------------------
// Setting Instance Pose
// ------------------------------------------------------------------------

func (self *SceneObjectPoses) ClearPoses() {
	self.poses = []geom2d.Matrix3{}
	self.poses_data_buffer = nil
	self.poses_webgl_buffer = js.Null()
}

func (self *SceneObjectPoses) AddPose(pose *geom2d.Matrix3) {
	self.poses = append(self.poses, *pose)
}

func (self *SceneObjectPoses) GetPoseCount() int {
	return len(self.poses)
}

// ------------------------------------------------------------------------
//  Data Buffers
// ------------------------------------------------------------------------

func (self *SceneObjectPoses) BuildDataBuffers() {
	// build data buffer for the instance poses
	self.poses_data_buffer = make([]float32, len(self.poses)*9)
	pos := 0
	for i := 0; i < len(self.poses); i++ {
		buf := self.poses_data_buffer[pos : pos+9]
		e := self.poses[i].GetElements()
		for j := 0; j < 9; j++ {
			buf[j] = e[j]
		}
		pos += 9
	}
	self.poses_webgl_buffer = js.Null()
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *SceneObjectPoses) IsWebGLBufferReady() bool {
	return !self.poses_webgl_buffer.IsNull()
}

func (self *SceneObjectPoses) build_webgl_buffers(wctx *common.WebGLContext) {
	// THIS FUCNTION IS MEANT TO BE CALLED BY RENDERER. NO NEED TO BE EXPORTED
	context := wctx.GetContext()     // js.Value
	constants := wctx.GetConstants() // *common.Constants
	if self.poses_data_buffer != nil {
		self.poses_webgl_buffer = context.Call("createBuffer", constants.ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, self.poses_webgl_buffer)
		var vertices_array = common.ConvertGoSliceToJsTypedArray(self.poses_data_buffer)
		context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
	} else {
		self.poses_webgl_buffer = js.Null()
	}
}

func (self *SceneObjectPoses) GetWebGLBuffer(mode string) (buffer js.Value, count int, pinfo [3]int) {
	if self.IsWebGLBufferReady() {
		return self.poses_webgl_buffer, len(self.poses_data_buffer), self.poses_info
	} else {
		return js.Null(), 0, [3]int{0, 0, 0}
	}
}
