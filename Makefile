all: 
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl_tester.go
	go build -o wasm_test_server wasm_test_server.go

2d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_tester.go
	go build -o wasm_test_server wasm_test_server.go

3d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl3d_tester.go
	go build -o wasm_test_server wasm_test_server.go

globe:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webglobe_tester.go
	go build -o wasm_test_server wasm_test_server.go

clean:
	rm wasm_test.wasm wasm_test_server
