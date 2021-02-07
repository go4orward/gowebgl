// This code launches a simple HTTP server to test the pre-built WASM bundle
// (This code is inspired by https://github.com/bobcob7/wasm-basic-triangle)
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

func main() {
	// load data for '/wasm_test.html'
	htmlData, err := ioutil.ReadFile("./wasm_test.html")
	if err != nil {
		log.Fatalf("Could not read wasm_test.html file: %s\n", err)
	}
	// load data for '/wasm/wasm_exec.js'
	wasmExecData, err := ioutil.ReadFile(runtime.GOROOT() + "/misc/wasm/wasm_exec.js")
	if err != nil {
		log.Fatalf("Could not read wasm_exec.js file: %s\n", err)
	}
	// load data for '/wasm/gowebgl.wasm'
	wasmData, err := ioutil.ReadFile("wasm_test.wasm")
	if err != nil {
		log.Fatalf("Could not read wasm file: %s\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(htmlData)
	})

	http.HandleFunc("/wasm/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(wasmExecData)
	})

	http.HandleFunc("/wasm/wasm_test.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		w.WriteHeader(http.StatusOK)
		w.Write(wasmData)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
