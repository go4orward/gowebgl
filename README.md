# GoWebGL : Interactive Graphics in Go

Interactive 2D & 3D Graphics Library using Go + WebAssembly + WebGL

## How to Build & Run

Simplest example: &emsp; _(for explaining how WebGL works)_
```bash
$ make                  # source : 'webgl_tester.go'
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webgl_teser result](assets/xscreen_webgl.png)

2D example: &emsp; _(with animation & user interactions)_
```bash
$ make 2d               # source : 'webgl2d_tester.go'
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
or
$ make 2dui             # source : 'webgl2d_ui_tester.go'
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webgl2d_teser result](assets/xscreen_webgl2d.png)

3D example: &emsp; _(with perspective & orthographic camera)_
```bash
$ make 3d               # source : 'webgl3d_tester.go'
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webgl3d_teser result](assets/xscreen_webgl3d.png)

Globe example: &emsp; _(with perspective & orthographic camera)_
```bash
$ make globe            # source : 'webglobe_tester.go'
$ ./wasm_test_server    # open your browser, and visit http://localhost:8080
```
![webglobe_teser result](assets/xscreen_webglobe.png)

## Thanks

I hope this project can help many people to learn WebGL and build awesome 2D & 3D graphics applications.  
Many thanks to [Richard Musiol](https://github.com/neelance), for his vision and contributions for GopherJS and WebAssembly for Go.

Resources from:
- [Go Gopher images](https://golang.org/doc/gopher/) originally created by Renee French
- [VisibleEarth by NASA](https://visibleearth.nasa.gov/collection/1484/blue-marble) : world satellite images
- [NaturalEarth](https://www.naturalearthdata.com/) : public domain map dataset
