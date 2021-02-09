package common

import (
	"errors"
	"fmt"
	"syscall/js"
)

type WebGLContext struct {
	width     int       //
	height    int       //
	canvas    js.Value  // canvas DOM element
	context   js.Value  // WebGL context object
	constants Constants // WebGL constant values
}

func NewWebGLContext(canvas_id string) (*WebGLContext, error) {
	var wctx WebGLContext
	// initialize the canvas
	doc := js.Global().Get("document")
	wctx.canvas = doc.Call("getElementById", canvas_id)
	if wctx.canvas.IsNull() {
		return nil, errors.New("Canvas not found (ID:'" + canvas_id + "')")
	}
	wctx.width = doc.Get("body").Get("clientWidth").Int()
	wctx.height = doc.Get("body").Get("clientHeight").Int()
	wctx.height = wctx.width // TODO: JUST FOR THE CASE WHERE THE LAST LINE FAILED, RETURNING 0
	wctx.canvas.Set("width", wctx.width)
	wctx.canvas.Set("height", wctx.height)

	// create WebGL context
	wctx.context = wctx.canvas.Call("getContext", "webgl")
	if wctx.context.IsUndefined() {
		wctx.context = wctx.canvas.Call("getContext", "experimental-webgl")
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

func (self *WebGLContext) GetCanvas() js.Value {
	return self.canvas
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

// ----------------------------------------------------------------------------
// User Interactions with Mouse
// ----------------------------------------------------------------------------

func (self *WebGLContext) SetEventHandler(handler_function_name string) {
	js.Global().Set(handler_function_name, EventHandlerWrapper())
}

func (self *WebGLContext) RegisterClickEventHandler(handler func(cxy [2]int, keystat [4]bool) bool) {
	evthandler_for_click = handler
}

func (self *WebGLContext) RegisterDoubleClickEventHandler(handler func(cxy [2]int, keystat [4]bool) bool) {
	evthandler_for_dblclick = handler
}

func (self *WebGLContext) RegisterMouseOverEventHandler(handler func(cxy [2]int, keystat [4]bool) bool) {
	evthandler_for_mouse_over = handler
}

func (self *WebGLContext) RegisterMouseDragEventHandler(handler func(cxy [2]int, sxy [2]int, keystat [4]bool) bool) {
	evthandler_for_mouse_drag = handler
}

func (self *WebGLContext) RegisterMouseWheelEventHandler(handler func(x int) bool) {
	evthandler_for_mouse_wheel = handler
}

var mouse_dragging bool = false
var mouse_sxy [2]int
var evthandler_for_click func(cxy [2]int, keystat [4]bool) bool = nil
var evthandler_for_dblclick func(cxy [2]int, keystat [4]bool) bool = nil
var evthandler_for_mouse_over func(cxy [2]int, keystat [4]bool) bool = nil
var evthandler_for_mouse_drag func(cxy [2]int, sxy [2]int, keystat [4]bool) bool = nil
var evthandler_for_mouse_wheel func(x int) bool = nil

func EventHandlerWrapper() js.Func {
	// NOTE THAT THIS WRAPPER FUNCTION SHOULD BE EXPORTED
	function := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		etype := event.Get("type").String()
		switch etype {
		case "click":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			dx, dy := cxy[0]-mouse_sxy[0], cxy[1]-mouse_sxy[1]
			if dx < -3 || dx > +3 || dy < -3 || dy > +3 {
				// ignore
			} else if evthandler_for_click != nil {
				evthandler_for_click(cxy, keystat)
			} else {
				fmt.Printf("%s %v %v\n", etype, cxy, keystat)
			}
		case "dblclick":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			if evthandler_for_dblclick != nil {
				evthandler_for_dblclick(cxy, keystat)
			} else {
				fmt.Printf("%s %v %v\n", etype, cxy, keystat)
			}
		case "mousemove":
			if mouse_dragging {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
				if evthandler_for_mouse_drag != nil {
					evthandler_for_mouse_drag(cxy, mouse_sxy, keystat)
				} else {
					fmt.Printf("%s %v => %v %v\n", etype, mouse_sxy, cxy, keystat)
				}
			} else {
				if evthandler_for_mouse_over != nil {
					cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
					keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
					evthandler_for_mouse_over(cxy, keystat)
				}
			}
		case "mousedown":
			mouse_sxy = [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			mouse_dragging = true
		case "mouseup":
			mouse_dragging = false
		case "mouseleave":
			mouse_dragging = false
		default:
			fmt.Println(etype)
		}
		return nil
	})
	return function
}
