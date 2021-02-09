package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_ForAxes(wctx *common.WebGLContext) *Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec2 pos;
		uniform mat3 mview;
		varying vec2 v_pos;
		void main() {
			v_pos = pos;
			vec3 xy1 = vec3(pos.x, pos.y, 1.0);
			vec3 new_pos = mview * xy1;
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec2 v_pos;
		void main() {
			if      (v_pos.x != 0.0) gl_FragColor = vec4(1.0, 0.2, 0.2, 1.0);
			else if (v_pos.y != 0.0) gl_FragColor = vec4(0.2, 1.0, 0.2, 1.0);
			else                     gl_FragColor = vec4(0.7, 0.7, 1.0, 1.0);
		}`
	shader, _ := NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("mview", "mat3", "renderer.modelview") // automatic binding of ModelView matrix of Renderer
	shader.InitBindingForAttribute("pos", "vec2", "geometry.coord")     // automatic binding of vertex coordinates
	return shader
}

func NewShader_SimplyRed(wctx *common.WebGLContext) *Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec2 pos;
		void main() {
			gl_Position = vec4(pos.x, pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		void main() {
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}`
	shader, _ := NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForAttribute("pos", "vec2", "") // 'SetBindingForAttribute()' has to be called later
	return shader
}

func NewShader_SingleColor(wctx *common.WebGLContext) *Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec2 pos;
		void main() {
			gl_Position = vec4(pos.x, pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("color", "vec3", "material.color") // automatic binding of uniform variable
	shader.InitBindingForAttribute("pos", "vec2", "geometry.coord") // automatic binding of attribute variable
	return shader
}

func NewShader_ModelView(wctx *common.WebGLContext) *Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec2 pos;
		uniform mat3 mview;
		void main() {
			vec3 xy1 = vec3(pos.x, pos.y, 1.0);
			vec3 new_pos = mview * xy1;
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("mview", "mat3", "renderer.modelview") // automatic binding of ModelView matrix of Renderer
	shader.InitBindingForUniform("color", "vec3", "material.color")     // automatic binding of uniform variable
	shader.InitBindingForAttribute("pos", "vec2", "geometry.coord")     // automatic binding of attribute variable
	return shader
}
