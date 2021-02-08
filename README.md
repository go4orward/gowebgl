# WebGL in Go

Interactive 2D & 3D Graphics Library using Go + WebAssembly + WebGL

## How to Build & Run

Simplest example: &emsp; _(for explaining how WebGL works)_
```bash
$ make 
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webgl_teser result](doc/xscreen_webgl.png)

2D example:
```bash
$ make 2d
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webgl2d_teser result](doc/xscreen_webgl2d.png)

3D example:
```bash
$ make 3d               # COMMING SOON 
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```

## Thanks

I hope this project to be helpful for many people to learn WebGL and build awesome 2D & 3D graphics applications.
Many thanks to @neelance (Richard Musiol), for his vision and huge contributions for GopherJS and WebAssembly for Go.
