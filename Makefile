all: 
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl1st_tester.go
	go build -o wasm_test_server wasm_test_server.go

dev:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_dev.go
	go build -o wasm_test_server wasm_test_server.go

2d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_tester.go
	go build -o wasm_test_server wasm_test_server.go

2dui:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl2d_ui_tester.go
	go build -o wasm_test_server wasm_test_server.go

3d:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webgl3d_tester.go
	go build -o wasm_test_server wasm_test_server.go

globe:
	GOOS=js GOARCH=wasm go build -o wasm_test.wasm webglobe_tester.go
	go build -o wasm_test_server wasm_test_server.go

clean:
	rm wasm_test.wasm wasm_test_server
