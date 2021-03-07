package wcommon

import (
	"fmt"
	"syscall/js"
)

type SceneObjectPoses struct {
	Size        int       //
	Count       int       //
	DataBuffer  []float32 //
	WebGLBuffer js.Value  //
}

func NewSceneObjectPoses(size int, count int, data []float32) *SceneObjectPoses {
	poses := SceneObjectPoses{Size: size, Count: count}
	poses.DataBuffer = make([]float32, size*count)
	if data != nil {
		for i := 0; i < len(poses.DataBuffer) && i < len(data); i++ {
			poses.DataBuffer[i] = data[i]
		}
	}
	poses.WebGLBuffer = js.Null()
	return &poses
}

func (self *SceneObjectPoses) ShowInfo() {
	fmt.Printf("Instance Poses : size = %d & count = %d\n", self.Size, self.Count)
}

// ------------------------------------------------------------------------
// Setting Instance Pose
// ------------------------------------------------------------------------

func (self *SceneObjectPoses) SetPose(index int, offset int, values ...float32) bool {
	if (offset + len(values)) > self.Size {
		return false
	}
	pos := index * self.Size
	for i := 0; i < len(values); i++ {
		self.DataBuffer[pos+offset+i] = values[i]
	}
	return true
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *SceneObjectPoses) IsWebGLBufferReady() bool {
	return !self.WebGLBuffer.IsNull()
}

func (self *SceneObjectPoses) BuildWebGLBuffer(wctx *WebGLContext) {
	// THIS FUCNTION IS MEANT TO BE CALLED BY RENDERER. NO NEED TO BE EXPORTED
	context := wctx.GetContext()     // js.Value
	constants := wctx.GetConstants() // *Constants
	if self.DataBuffer != nil {
		self.WebGLBuffer = context.Call("createBuffer", constants.ARRAY_BUFFER)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, self.WebGLBuffer)
		var vertices_array = ConvertGoSliceToJsTypedArray(self.DataBuffer)
		context.Call("bufferData", constants.ARRAY_BUFFER, vertices_array, constants.STATIC_DRAW)
		context.Call("bindBuffer", constants.ARRAY_BUFFER, nil)
	} else {
		self.WebGLBuffer = js.Null()
	}
}
