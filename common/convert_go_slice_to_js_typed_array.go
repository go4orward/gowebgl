package common

import (
	"fmt"
	"reflect"
	"strconv"
	"syscall/js"
	"unsafe"
)

// Since js.TypedArrayOf() of Go1.11 is no longer supported (due to WASM memory issue),
// we have to use js.CopyBytesToJS() instead. (Now it runs fine with Go1.15.7, Feb 5 2021)
//   Ref: syscall/js: replace TypedArrayOf with CopyBytesToGo/CopyBytesToJS
//   Ref: https://github.com/golang/go/issues/31980  	("js.TypedArrayOf is impossible to use correctly")
//   Ref: https://go-review.googlesource.com/c/go/+/177537/
//   Ref: https://github.com/golang/go/issues/32402  	(solution provided by 'hajimehoshi')
//   Ref: https://github.com/nuberu/webgl				(Golang WebAssembly wrapper for WebGL)
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
		// ShowArrayInfo(byte_slice)
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, byte_slice)
		// ShowArrayInfo(a)
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

func ShowArrayInfo(desc string, iv interface{}) {
	switch iv.(type) {
	case []byte:
		v := iv.([]byte)
		fmt.Printf("ARRAY %s of []byte : len:%d  %v\n", desc, len(v), v)
	case []uint32:
		v := iv.([]uint32)
		fmt.Printf("ARRAY %s of []uint32 : len:%d  %v\n", desc, len(v), v)
	case []float32:
		v := iv.([]float32)
		fmt.Printf("ARRAY %s of []float32 : len:%d  %v\n", desc, len(v), v)
	case js.Value:
		v := iv.(js.Value)
		if v.IsUndefined() {
			fmt.Printf("ARRAY %s of js.Value : undefined\n", desc)
		} else if v.IsNull() {
			fmt.Printf("ARRAY %s of js.Value : null\n", desc)
		} else {
			fmt.Printf("ARRAY %s of js.Value : type:%s name:%s length:%s (%s) byteLength:%s byteOffset:%s\n", desc, v.Type().String(), v.Get("name").String(), p_int(v, "length"), p_int(v, "BYTES_PER_ELEMENT"), p_int(v, "byteLength"), p_int(v, "byteOffset"))
			fmt.Printf("ARRAY %s of js.Value : len:%d  %v\n", desc, v.Length(), contents(v))
		}
	default:
		fmt.Printf("ARRAY %s of UNKNOWN \n", desc)
	}
}

func p_int(v js.Value, property string) string {
	p := v.Get(property)
	if p.IsUndefined() {
		return "undefined"
	} else {
		return strconv.Itoa(p.Int())
	}
}
func contents(v js.Value) string {
	contents := ""
	for i := 0; i < v.Length(); i++ {
		if i == 0 {
			contents += v.Index(i).String()
		} else {
			contents += "," + v.Index(i).String()
		}
	}
	return "[" + contents + "]"
}
