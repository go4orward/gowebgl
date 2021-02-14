package common

import (
	"syscall/js"
)

type Constants struct {
	ARRAY_BUFFER         js.Value //
	COLOR_BUFFER_BIT     js.Value //
	COMPILE_STATUS       js.Value //
	DEPTH_BUFFER_BIT     js.Value //
	DEPTH_TEST           js.Value //
	ELEMENT_ARRAY_BUFFER js.Value //
	FLOAT                js.Value //
	FRAGMENT_SHADER      js.Value //
	LEQUAL               js.Value //
	LINES                js.Value //
	LINK_STATUS          js.Value //
	POINTS               js.Value //
	STATIC_DRAW          js.Value //
	TRIANGLES            js.Value //
	UNSIGNED_INT         js.Value //
	UNSIGNED_SHORT       js.Value //
	VERTEX_SHADER        js.Value //
}

func (self *Constants) LoadFromContext(context js.Value) {
	// get WebGL constants
	self.ARRAY_BUFFER = context.Get("ARRAY_BUFFER")
	self.COLOR_BUFFER_BIT = context.Get("COLOR_BUFFER_BIT")
	self.COMPILE_STATUS = context.Get("COMPILE_STATUS")
	self.DEPTH_BUFFER_BIT = context.Get("DEPTH_BUFFER_BIT")
	self.DEPTH_TEST = context.Get("DEPTH_TEST")
	self.ELEMENT_ARRAY_BUFFER = context.Get("ELEMENT_ARRAY_BUFFER")
	self.FLOAT = context.Get("FLOAT")
	self.FRAGMENT_SHADER = context.Get("FRAGMENT_SHADER")
	self.LEQUAL = context.Get("LEQUAL")
	self.LINES = context.Get("LINES")
	self.LINK_STATUS = context.Get("LINK_STATUS")
	self.POINTS = context.Get("POINTS")
	self.STATIC_DRAW = context.Get("STATIC_DRAW")
	self.TRIANGLES = context.Get("TRIANGLES")
	self.UNSIGNED_INT = context.Get("UNSIGNED_INT")
	self.UNSIGNED_SHORT = context.Get("UNSIGNED_SHORT")
	self.VERTEX_SHADER = context.Get("VERTEX_SHADER")
}
