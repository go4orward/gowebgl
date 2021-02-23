all: 
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl1st_example.go
	go build -o wasm_test_server wasm_test_server.go

2d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_example.go
	go build -o wasm_test_server wasm_test_server.go

2dui:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_ui_example.go
	go build -o wasm_test_server wasm_test_server.go

3d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl3d_example.go
	go build -o wasm_test_server wasm_test_server.go

globe:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webglobe_example.go
	go build -o wasm_test_server wasm_test_server.go

clean:
	rm wasm_test.wasm wasm_test_server
