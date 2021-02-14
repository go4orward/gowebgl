package common

import (
	"errors"
	"fmt"
	"math"
	"syscall/js"
)

type WebGLContext struct {
	width     int       //
	height    int       //
	canvas_id string    // canvas DOM element's ID
	canvas    js.Value  // canvas DOM element
	context   js.Value  // WebGL context object
	constants Constants // WebGL constant values
	ext_uint  js.Value  // extension for "OES_element_index_uint"
	ext_angle js.Value  // extension for "ANGLE_instanced_arrays"
}

func NewWebGLContext(canvas_id string) (*WebGLContext, error) {
	var wctx WebGLContext
	// initialize the canvas
	doc := js.Global().Get("document")
	wctx.canvas = doc.Call("getElementById", canvas_id)
	if wctx.canvas.IsNull() {
		return nil, errors.New("Canvas not found (ID:'" + canvas_id + "')")
	}
	wctx.canvas_id = canvas_id
	if true {
		wctx.width = wctx.canvas.Get("clientWidth").Int()
		wctx.height = wctx.canvas.Get("clientHeight").Int()
	} else {
		wctx.width = doc.Get("body").Get("clientWidth").Int()
		wctx.height = doc.Get("body").Get("clientHeight").Int()
	}
	wctx.canvas.Set("width", wctx.width)   // IMPORTANT!
	wctx.canvas.Set("height", wctx.height) // IMPORTANT!
	// create WebGL context
	wctx.context = wctx.canvas.Call("getContext", "webgl")
	if wctx.context.IsUndefined() {
		wctx.context = wctx.canvas.Call("getContext", "experimental-webgl")
		if wctx.context.IsUndefined() {
			return nil, errors.New("WebGL not supported")
		}
	}
	wctx.constants.LoadFromContext(wctx.context) // load WebGL constants
	wctx.SetupExtension("UINT")                  // extension for UINT32 index
	// wctx.SetupExtension("ANGLE")              // extension for geometry instancing (Renderer will set it up)
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

func (self *WebGLContext) GetWH() [2]int {
	// AspectRatio : "width : height", "width to height", or "width / height"
	return [2]int{self.width, self.height}
}

func (self *WebGLContext) ShowInfo() {
	fmt.Printf("WebGLContext : canvas '%s' (%d x %d)\n", self.canvas_id, self.width, self.height)
}

// ----------------------------------------------------------------------------
// WebGL Extensions
// ----------------------------------------------------------------------------

func (self *WebGLContext) SetupExtension(extname string) {
	switch extname {
	case "UINT": // extension for UINT32 index, to drawElements() with large number of vertices
		self.ext_uint = self.context.Call("getExtension", "OES_element_index_uint")
	case "ANGLE": // extension for geometry instancing
		self.ext_angle = self.context.Call("getExtension", "ANGLE_instanced_arrays")
	}
}

func (self *WebGLContext) GetExtension(extname string) js.Value {
	switch extname {
	case "UINT": // extension for UINT32 index, to drawElements() with large number of vertices
		return self.ext_uint
	case "ANGLE": // extension for geometry instancing
		return self.ext_angle
	}
	return js.Null()
}

func (self *WebGLContext) IsExtensionReady(extname string) bool {
	switch extname {
	case "UINT": // extension for UINT32 index, to drawElements() with large number of vertices
		return !self.ext_uint.IsNull() && !self.ext_uint.IsUndefined()
	case "ANGLE": // extension for geometry instancing
		return !self.ext_angle.IsNull() && !self.ext_angle.IsUndefined()
	}
	return false
}

// ----------------------------------------------------------------------------
// User Interactions (Event Handling)
// ----------------------------------------------------------------------------

func (self *WebGLContext) SetupEventHandlers() {
	// export EventHandling function from Go side
	js.Global().Set("goEventHandler", go_wrapper_for_event_handler())
	// add EventListener functions from Javascript side
	wasm_js_listener := js.Global().Get("wasm_js_listener") // 'wasm_js_listener()' must call 'goEventHandler()'
	if wasm_js_listener.IsUndefined() {
		fmt.Println("Setting up EventHandler failed : 'wasm_js_listener' function not found")
		fmt.Println("  (for example, 'wasm_js_listener = function(event){goEventHandler(event);};')")
	} else {
		self.canvas.Call("addEventListener", "click", wasm_js_listener)
		self.canvas.Call("addEventListener", "dblclick", wasm_js_listener)
		self.canvas.Call("addEventListener", "mousemove", wasm_js_listener)
		self.canvas.Call("addEventListener", "mousedown", wasm_js_listener)
		self.canvas.Call("addEventListener", "mouseup", wasm_js_listener)
		self.canvas.Call("addEventListener", "mouseleave", wasm_js_listener)
		self.canvas.Call("addEventListener", "wheel", wasm_js_listener)
		js.Global().Get("window").Call("addEventListener", "resize", wasm_js_listener)
		// What it actually does is like:
		// canvas.addEventListener("click", function(event) { goEventHandler(canvas, event); });
	}
}

func (self *WebGLContext) RegisterEventHandlerForClick(handler func(canvasxy [2]int, keystat [4]bool)) {
	evthandler_for_click = handler
}

func (self *WebGLContext) RegisterEventHandlerForDoubleClick(handler func(canvasxy [2]int, keystat [4]bool)) {
	evthandler_for_dblclick = handler
}

func (self *WebGLContext) RegisterEventHandlerForMouseOver(handler func(canvasxy [2]int, keystat [4]bool)) {
	evthandler_for_mouse_over = handler
}

func (self *WebGLContext) RegisterEventHandlerForMouseDrag(handler func(canvasxy [2]int, dxy [2]int, keystat [4]bool)) {
	evthandler_for_mouse_drag = handler // 'dx' & 'dy' is delta movement in Camera space coordinates
}

func (self *WebGLContext) RegisterEventHandlerForMouseWheel(handler func(canvasxy [2]int, scale float32, keystat [4]bool)) {
	evthandler_for_mouse_wheel = handler // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
}

func (self *WebGLContext) RegisterEventHandlerForWindowResize(handler func(w int, h int)) {
	evthandler_for_window_resize = handler
}

var mouse_dragging bool = false
var mouse_sxy = [2]int{0, 0}
var mouse_wheel_scale float64 = 500 // in the range of [0 ~ 500(default) ~ 1000]
var evthandler_for_click func(canvasxy [2]int, keystat [4]bool) = nil
var evthandler_for_dblclick func(canvasxy [2]int, keystat [4]bool) = nil
var evthandler_for_mouse_over func(canvasxy [2]int, keystat [4]bool) = nil
var evthandler_for_mouse_drag func(canvasxy [2]int, dxy [2]int, keystat [4]bool) = nil
var evthandler_for_mouse_wheel func(canvasxy [2]int, scale float32, keystat [4]bool) = nil // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
var evthandler_for_window_resize func(window_width int, window_height int) = nil

func go_wrapper_for_event_handler() js.Func {
	// NOTE THAT THIS WRAPPER FUNCTION SHOULD BE EXPORTED
	function := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			fmt.Println("Invalid GoCallback call (for EventHandling) from Javascript")
			return nil
		}
		event := args[0]                    // js.Value (event object)
		etype := event.Get("type").String() // canvas := event.Get("srcElement")
		switch etype {
		case "click":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			dx, dy := (cxy[0] - mouse_sxy[0]), (cxy[1] - mouse_sxy[1])
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			if dx < -3 || dx > +3 || dy < -3 || dy > +3 {
				// ignore
			} else if evthandler_for_click != nil {
				evthandler_for_click(cxy, keystat)
			} else {
				fmt.Printf("%s (%d %d) %v\n", etype, cxy[0], cxy[1], keystat)
			}
		case "dblclick":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			if evthandler_for_dblclick != nil {
				evthandler_for_dblclick(cxy, keystat)
			} else {
				fmt.Printf("%s (%d %d) %v\n", etype, cxy[0], cxy[1], keystat)
			}
		case "mousemove":
			if mouse_dragging {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				dxy := [2]int{event.Get("movementX").Int(), event.Get("movementY").Int()}
				keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
				if evthandler_for_mouse_drag != nil {
					evthandler_for_mouse_drag(cxy, dxy, keystat)
				} else {
					fmt.Printf("%s (dx dy) with %v\n", etype, dxy[0], dxy[1], keystat)
				}
			} else {
				if evthandler_for_mouse_over != nil {
					cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
					keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
					evthandler_for_mouse_over(cxy, keystat)
				}
			}
		case "mousedown":
			mouse_dragging = true
			mouse_sxy = [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
		case "mouseup":
			mouse_dragging = false
		case "mouseleave":
			mouse_dragging = false
		case "wheel":
			if evthandler_for_mouse_wheel != nil {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				mouse_wheel_scale += float64(event.Get("deltaY").Int()) // [ 0 ~ 500(default) ~ 1000 ]
				mouse_wheel_scale = float64(math.Max(0, math.Min(mouse_wheel_scale, 1000)))
				scale_exp := (mouse_wheel_scale - 500.0) / 250.0 // [ -2 ~ 0(default) ~ +2 ]
				scale := math.Pow(10, scale_exp)                 // [ 0.01 ~ 1(default) ~ 100.0 ]
				keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
				evthandler_for_mouse_wheel(cxy, float32(scale), keystat)
			}
		case "resize":
			w := js.Global().Get("window").Get("innerWidth").Int()
			h := js.Global().Get("window").Get("innerHeight").Int()
			if evthandler_for_window_resize != nil {
				evthandler_for_window_resize(w, h)
			} else {
				fmt.Printf("window.resize %d %d\n", w, h)
			}
		default:
			fmt.Println(etype)
		}
		return nil
	})
	return function
}

// ----------------------------------------------------------------------------
// Animation Frame
// ----------------------------------------------------------------------------

func (self *WebGLContext) SetupAnimationFrame() {
	// export EventHandling function from Go side
	js.Global().Set("goSceneRenderer", go_wrapper_for_animation_frame())
	// add EventListener functions from Javascript side
	wasm_js_renderer := js.Global().Get("wasm_js_renderer") // 'wasm_js_renderer()' must call 'goSceneRenderer()'
	if wasm_js_renderer.IsUndefined() {
		fmt.Println("Setting up EventHandler failed : 'wasm_js_renderer' function not found")
		fmt.Println("  (for example, 'wasm_js_renderer = function(){goSceneRenderer();}')")
	} else {
		// What it actually does is like:
		//   requestAnimationFrame(drawSceneForAnimation);
		//   function drawSceneForAnimation() {
		//     if (typeof goDrawAnimationFrame != 'undefined') {
		//         goDrawAnimationFrame(canvas);   // draw the scene by calling Go renderer function
		//     }
		//     requestAnimationFrame(drawSceneForAnimation); // call itself again for the next frame
		//   }
		js.Global().Call("requestAnimationFrame", wasm_js_renderer, self.canvas)
	}
}

// RegisterDrawSceneCallback
func (self *WebGLContext) RegisterDrawHandlerForAnimationFrame(function func(canvas js.Value)) {
	handler_draw_animation_frame = function
}

var handler_draw_animation_frame func(canvas js.Value) = nil

func go_wrapper_for_animation_frame() js.Func {
	// NOTE THAT THIS WRAPPER FUNCTION SHOULD BE EXPORTED
	function := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			fmt.Println("Invalid GoCallback call (for AnimationFrame) from Javascript")
			return nil
		}
		canvas := args[0] // js.Value (canvas DOM element)
		if handler_draw_animation_frame != nil {
			handler_draw_animation_frame(canvas)
		}
		wasm_js_renderer := js.Global().Get("wasm_js_renderer")
		js.Global().Call("requestAnimationFrame", wasm_js_renderer, canvas)
		return nil
	})
	return function
}
