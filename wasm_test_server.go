// This code launches a simple HTTP server to test the pre-built WASM bundle
// (This code was inspired by https://github.com/bobcob7/wasm-basic-triangle)
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

func main() {
	// HTML file
	html, err := ioutil.ReadFile("./wasm_test.html")
	if err != nil {
		log.Fatalf("Could not read wasm_test.html file: %s\n", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
	})

	// texture images
	http.Handle("/assets/", http.FileServer(http.Dir(".")))

	// 'wasm_exec.js'
	exjs, err := ioutil.ReadFile(runtime.GOROOT() + "/misc/wasm/wasm_exec.js")
	if err != nil {
		log.Fatalf("Could not read wasm_exec.js file: %s\n", err)
	}
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(exjs)
	})

	// 'wasm_test.wasm'
	wasm, err := ioutil.ReadFile("wasm_test.wasm")
	if err != nil {
		log.Fatalf("Could not read wasm file: %s\n", err)
	}
	http.HandleFunc("/wasm_test.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		w.WriteHeader(http.StatusOK)
		w.Write(wasm)
	})

	// start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
