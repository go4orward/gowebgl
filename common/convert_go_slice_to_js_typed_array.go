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

func ConvertGoSliceToJsTypedArray(a interface{}) js.Value {
	switch a := a.(type) {
	case []int8:
		b := js.Global().Get("Uint8Array").New(len(a))
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int8Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength"))
	case []int16:
		b := js.Global().Get("Uint8Array").New(len(a) * 2)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int16Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/2)
	case []int32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []int64:
		b := js.Global().Get("Uint8Array").New(len(a) * 8)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("BigInt64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	case []uint8:
		b := js.Global().Get("Uint8Array").New(len(a))
		js.CopyBytesToJS(b, a)
		return b
	case []uint16:
		b := js.Global().Get("Uint8Array").New(len(a) * 2)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Uint16Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/2)
	case []uint32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Uint32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []uint64:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("BigUint64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	case []float32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		// ShowArrayInfo(byte_slice)
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Float32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []float64:
		b := js.Global().Get("Uint8Array").New(len(a) * 8)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Float64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("Unexpected value at ConvertGoSliceToJsTypedArray(): %T", a))
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
			fmt.Printf("ARRAY %s of js.Value : type:%s name:%s (%s) byteLength:%s length:%s byteOffset:%s\n", desc, v.Type().String(), v.Get("name").String(), p_int(v, "length"), p_int(v, "byteLength"), p_int(v, "BYTES_PER_ELEMENT"), p_int(v, "byteOffset"))
			fmt.Printf("ARRAY %s of js.Value : len:%d  %v\n", desc, v.Length(), a_contents(v))
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

func a_contents(v js.Value) string {
	contents := ""
	for i := 0; i < v.Length(); i++ {
		if i == 0 { // v.Index(i).Type().String() == "number"
			contents += fmt.Sprintf("%.1f", v.Index(i).Float())
		} else {
			contents += "," + fmt.Sprintf("%.1f", v.Index(i).Float())
		}
	}
	return "[" + contents + "]"
}
