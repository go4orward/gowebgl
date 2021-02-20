package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	ui1, ui3 := uint16(1), uint16(3)
	ui13 := uint32(ui1)<<16 + uint32(ui3)
	fmt.Printf("(1) %#04x  (3) %#04x  (1<<16+3) %#08x\n", ui1, ui3, ui13)

	b1, b3, b13 := []byte{0, 0}, []byte{0, 0}, []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b1, ui1)
	binary.LittleEndian.PutUint16(b3, ui3)
	binary.LittleEndian.PutUint32(b13, ui13)

	fmt.Printf("(1) %v  (3) %v  (1<<16+3) %v\n", b1, b3, b13)

	// if true { // FOR DEBUGGING ONLY
	// 	ui, vi := uint32(65535*u), uint32(65535*v)
	// 	uib, vib := strconv.FormatUint(uint64(ui), 2), strconv.FormatUint(uint64(vi), 2)
	// 	final := math.Float32frombits(uint32(65535*u) + uint32(65535*v)<<16)
	// 	fi := math.Float32bits(final)
	// 	fib := strconv.FormatUint(uint64(fi), 2)
	// 	fmt.Printf("uv=(%.2f %.2f) u=%d (%032v) v=%d (%032v)  fi=%d (%032v)\n", u, v, ui, uib, vi, vib, fi, fib)
	// }
}
