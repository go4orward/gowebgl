package common

import (
	"fmt"
	"reflect"
	"syscall/js"
	"unsafe"
)

// Since js.TypedArrayOf() of Go1.11 is no longer supported (due to WASM memory issue),
// we have to use js.CopyBytesToJS() instead. (Now it runs fine with Go1.15.7, Feb 5 2021)
//   Ref: syscall/js: replace TypedArrayOf with CopyBytesToGo/CopyBytesToJS
//   Ref: https://github.com/golang/go/issues/31980  ("js.TypedArrayOf is impossible to use correctly")
//   Ref: https://go-review.googlesource.com/c/go/+/177537/
//   Ref: https://github.com/golang/go/issues/32402  (solution provided by 'hajimehoshi')
// Note that this solution sacrifices performance. (WebGL renderer's frame rate will be OK, though)
// We hope Go/WebAssembly will sort out this issue in the future.

func ConvertGoSliceToJsTypedArray(s interface{}) js.Value {
	switch s := s.(type) {
	case []int8:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Int8Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength"))
	case []int16:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Int16Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []int32:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Int32Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []int64:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 8)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("BigInt64Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	case []uint8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, s)
		return a
	case []uint16:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Uint16Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []uint32:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Uint32Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []uint64:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("BigUint64Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	case []float32:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Float32Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float64:
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		a := js.Global().Get("Uint8Array").New(len(s) * 8)
		js.CopyBytesToJS(a, byte_slice)
		return js.Global().Get("Float64Array").New(a.Get("buffer"), a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("Unexpected value at ConvertGoSliceToJsTypedArray(): %T", s))
	}
}
