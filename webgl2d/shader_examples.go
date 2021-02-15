package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_2DAxes(wctx *common.WebGLContext) *common.Shader {
	// Shader for three axes - X(RED) & Y(GREEN) - for visual reference,
	//   with auto-binded (Proj * View * Model) matrix and XY coordinates
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;
		attribute vec2 xy;
		varying   vec2 vxy;
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			vxy = xy;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec2 vxy;
		void main() {
			if (vxy.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else              gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")     // automatic binding of Proj*View*Model matrix
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords") // automatic binding of point coordinates
	return shader
}

func NewShader_SimplyRed(wctx *common.WebGLContext) *common.Shader {
	// Shader with with XY coordinates only (assuming Color == RED and PVM == Identity)
	var vertex_shader_code = `
		precision mediump float;
		attribute vec2 xy;
		void main() {
			gl_Position = vec4(xy.x, xy.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		void main() {
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords") // automatic binding of point coordinates
	return shader
}

func NewShader_Basic(wctx *common.WebGLContext) *common.Shader {
	// Shader with auto-binded color and (Proj * View * Model) matrix
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;
		attribute vec2 xy;
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")     // automatic binding of Proj*View*Model matrix
	shader.InitBindingForUniform("color", "vec3", "material.color") // automatic binding of material color
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords") // automatic binding of point coordinates
	return shader
}

func NewShader_InstancePoseColor(wctx *common.WebGLContext) *common.Shader {
	// Shader with instance pose, for rendering multiple instances of a same geometry
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;
		attribute vec2 xy;
		attribute vec2 txy;
		attribute vec3 color;
		varying   vec3 vc;
		void main() {
			vec3 new_pos = pvm * vec3(xy.x + txy.x, xy.y + txy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			vc = color;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying   vec3 vc;
		void main() { 
			gl_FragColor = vec4(vc.r, vc.g, vc.b, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")          // automatic binding of Proj*View*Model matrix
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords")      // automatic binding of point coordinates
	shader.InitBindingForAttribute("txy", "vec2", "instance.pose:5:0")   // automatic binding of instance pose
	shader.InitBindingForAttribute("color", "vec3", "instance.pose:5:2") // automatic binding of instance color
	return shader
}
