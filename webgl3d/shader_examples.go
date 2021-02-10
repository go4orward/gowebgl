package webgl3d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_ForAxes(wctx *common.WebGLContext) *common.Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec3 pos;
		uniform mat4 mview;
		varying vec3 v_pos;
		void main() {
			v_pos = pos;
			vec4 xy1 = vec4(pos.x, pos.y, pos.z, 1.0);
			vec4 new_pos = mview * xy1;
			gl_Position = vec4(new_pos.x, new_pos.y, new_pos.z, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_pos;
		void main() {
			if      (v_pos.x != 0.0) gl_FragColor = vec4(1.0, 0.2, 0.2, 1.0);
			else if (v_pos.y != 0.0) gl_FragColor = vec4(0.2, 1.0, 0.2, 1.0);
			else                     gl_FragColor = vec4(0.7, 0.7, 1.0, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("mview", "mat4", "renderer.modelview") // automatic binding of ModelView matrix of Renderer
	shader.InitBindingForAttribute("pos", "vec3", "geometry.coord:0:0") // automatic binding of vertex coordinates
	return shader
}

func NewShader_Basic(wctx *common.WebGLContext) *common.Shader {
	var vertex_shader_code = `
		precision mediump float;
		attribute vec3 pos;
		uniform mat4 mview;
		void main() {
			vec4 xy1 = vec4(pos.x, pos.y, pos.z, 1.0);
			vec4 new_pos = mview * xy1;
			gl_Position = vec4(new_pos.x, new_pos.y, new_pos.z, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("mview", "mat4", "renderer.modelview") // automatic binding of ModelView matrix of Renderer
	shader.InitBindingForUniform("color", "vec3", "material.color")     // automatic binding of uniform variable
	shader.InitBindingForAttribute("pos", "vec3", "geometry.coord:0:0") // automatic binding of attribute variable
	return shader
}
