package wcommon

import "syscall/js"

type Geometry interface {
	IsDataBufferReady() bool
	IsWebGLBufferReady() bool
	BuildWebGLBuffers(wctx *WebGLContext, for_points bool, for_lines bool, for_faces bool)
	GetWebGLBuffer(draw_mode int) (js.Value, int, [4]int)
	ShowInfo()
}
