package wcommon

import (
	"syscall/js"
)

type Constants struct {
	ARRAY_BUFFER         js.Value //
	BLEND                js.Value // for gl.enable(gl.BLEND)
	BYTE                 js.Value //
	CLAMP_TO_EDGE        js.Value // for gl.texParameteri()
	COLOR_BUFFER_BIT     js.Value //
	COMPILE_STATUS       js.Value //
	DEPTH_BUFFER_BIT     js.Value //
	DEPTH_TEST           js.Value //
	ELEMENT_ARRAY_BUFFER js.Value //
	FLOAT                js.Value //
	FRAGMENT_SHADER      js.Value //
	LEQUAL               js.Value //
	LINEAR               js.Value // for gl.texParameteri()
	LINES                js.Value //
	LINK_STATUS          js.Value //
	NEAREST              js.Value // for gl.texParameteri()
	ONE                  js.Value // for gl.blendFunc()
	ONE_MINUS_SRC_ALPHA  js.Value // for gl.blendFunc()
	POINTS               js.Value //
	RGBA                 js.Value //
	SRC_ALPHA            js.Value // for gl.blendFunc()
	STATIC_DRAW          js.Value //
	TEXTURE_2D           js.Value // for gl.texParameteri()
	TEXTURE0             js.Value //
	TEXTURE1             js.Value //
	TEXTURE_MIN_FILTER   js.Value // for gl.texParameteri()
	TEXTURE_WRAP_S       js.Value // for gl.texParameteri()
	TEXTURE_WRAP_T       js.Value // for gl.texParameteri()
	TRIANGLES            js.Value //
	UNSIGNED_BYTE        js.Value //
	UNSIGNED_INT         js.Value //
	UNSIGNED_SHORT       js.Value //
	VERTEX_SHADER        js.Value //
}

func (self *Constants) LoadFromContext(context js.Value) {
	// get WebGL constants
	self.ARRAY_BUFFER = context.Get("ARRAY_BUFFER")
	self.BLEND = context.Get("BLEND")
	self.BYTE = context.Get("BYTE")
	self.CLAMP_TO_EDGE = context.Get("CLAMP_TO_EDGE")
	self.COLOR_BUFFER_BIT = context.Get("COLOR_BUFFER_BIT")
	self.COMPILE_STATUS = context.Get("COMPILE_STATUS")
	self.DEPTH_BUFFER_BIT = context.Get("DEPTH_BUFFER_BIT")
	self.DEPTH_TEST = context.Get("DEPTH_TEST")
	self.ELEMENT_ARRAY_BUFFER = context.Get("ELEMENT_ARRAY_BUFFER")
	self.FLOAT = context.Get("FLOAT")
	self.FRAGMENT_SHADER = context.Get("FRAGMENT_SHADER")
	self.LEQUAL = context.Get("LEQUAL")
	self.LINEAR = context.Get("LINEAR")
	self.LINES = context.Get("LINES")
	self.LINK_STATUS = context.Get("LINK_STATUS")
	self.NEAREST = context.Get("NEAREST")
	self.ONE = context.Get("ONE")
	self.ONE_MINUS_SRC_ALPHA = context.Get("ONE_MINUS_SRC_ALPHA")
	self.POINTS = context.Get("POINTS")
	self.RGBA = context.Get("RGBA")
	self.SRC_ALPHA = context.Get("SRC_ALPHA")
	self.STATIC_DRAW = context.Get("STATIC_DRAW")
	self.TEXTURE_2D = context.Get("TEXTURE_2D")
	self.TEXTURE0 = context.Get("TEXTURE0")
	self.TEXTURE1 = context.Get("TEXTURE1")
	self.TEXTURE_MIN_FILTER = context.Get("TEXTURE_MIN_FILTER")
	self.TEXTURE_WRAP_S = context.Get("TEXTURE_WRAP_S")
	self.TEXTURE_WRAP_T = context.Get("TEXTURE_WRAP_T")
	self.TRIANGLES = context.Get("TRIANGLES")
	self.UNSIGNED_BYTE = context.Get("UNSIGNED_BYTE")
	self.UNSIGNED_INT = context.Get("UNSIGNED_INT")
	self.UNSIGNED_SHORT = context.Get("UNSIGNED_SHORT")
	self.VERTEX_SHADER = context.Get("VERTEX_SHADER")
}
