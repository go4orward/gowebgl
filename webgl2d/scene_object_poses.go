package webgl2d

import (
	"fmt"
	"syscall/js"

	"github.com/go4orward/gowebgl/common"
)

type SceneObjectPoses struct {
	size         int       //
	count        int       //
	data_buffer  []float32 //
	webgl_buffer js.Value  //
}

func NewSceneObjectPoses(size int, count int, data []float32) *SceneObjectPoses {
	poses := SceneObjectPoses{size: size, count: count}
	poses.data_buffer = make([]float32, size*count)
	if data != nil {
		for i := 0; i < len(poses.data_buffer) && i < len(data); i++ {
			poses.data_buffer[i] = data[i]
		}
	}
	poses.webgl_buffer = js.Null()
	return &poses
}

func (self *SceneObjectPoses) GetPoseSize() int {
	return self.size
}

func (self *SceneObjectPoses) GetPoseCount() int {
	return self.count
}

func (self *SceneObjectPoses) ShowInfo() {
	fmt.Printf("Instance Poses : size = %d & count = %d\n", self.size, self.count)
}

// ------------------------------------------------------------------------
// Setting Instance Pose
// ------------------------------------------------------------------------

func (self *SceneObjectPoses) SetPose(index int, offset int, values ...float32) bool {
	if (offset + len(values)) > self.size {
		return false
	}
	pos := index * self.size
	for i := 0; i < len(values); i++ {
		self.data_buffer[pos+offset+i] = values[i]
	}
	return true
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *SceneObjectPoses) IsWebGLBufferReady() bool {
	return !self.webgl_buffer.IsNull()
}

func (self *SceneObjectPoses) BuildWebGLBuffers(wctx *common.WebGLContext) {
	// THIS FUCNTION IS MEANT TO BE CALLED BY RENDERER. NO NEED TO BE EXPORTED
	context := wctx.GetContext()     // js.Value
	constants := wctx.GetConstants() // *common.Constants
	if self.data_buffer != nil {
		self.webgl_buffer = context.Call("createBuffer", constants.ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, self.webgl_buffer)
		var vertices_array = common.ConvertGoSliceToJsTypedArray(self.data_buffer)
		context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
	} else {
		self.webgl_buffer = js.Null()
	}
}
