// Package webgl3d implements the 3D Graphics library.
package common

import (
	"errors"
	"syscall/js"
)

type WebGLContext struct {
	width     int       //
	height    int       //
	context   js.Value  // WebGL context object
	constants Constants // WebGL constant values
}

func NewWebGLContext(canvas_id string) (*WebGLContext, error) {
	var wctx WebGLContext
	// initialize the canvas
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", canvas_id)
	if canvasEl.IsNull() {
		return nil, errors.New("Canvas not found (ID:'" + canvas_id + "')")
	}
	wctx.width = doc.Get("body").Get("clientWidth").Int()
	wctx.height = doc.Get("body").Get("clientHeight").Int()
	wctx.height = wctx.width // TODO: JUST FOR THE CASE WHERE THE LAST LINE FAILED, RETURNING 0
	canvasEl.Set("width", wctx.width)
	canvasEl.Set("height", wctx.height)

	// create WebGL context
	wctx.context = canvasEl.Call("getContext", "webgl")
	if wctx.context.IsUndefined() {
		wctx.context = canvasEl.Call("getContext", "experimental-webgl")
		if wctx.context.IsUndefined() {
			return nil, errors.New("WebGL not supported")
		}
	}
	wctx.constants.LoadFromContext(wctx.context) // load WebGL constants

	// Get Extension for UINT index for drawElements() with large number of vertices
	wctx.context.Call("getExtension", "OES_element_index_uint")
	// Get extension for geometry instancing
	// wctx.context.Call("getExtension", "ANGLE_instanced_arrays")
	return &wctx, nil
}

func (self *WebGLContext) GetContext() js.Value {
	return self.context
}

func (self *WebGLContext) GetConstants() *Constants {
	return &self.constants
}

func (self *WebGLContext) GetWidth() int {
	return self.width
}

func (self *WebGLContext) GetHeight() int {
	return self.height
}
