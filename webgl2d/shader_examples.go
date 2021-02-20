package webgl2d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_2DAxes(wctx *common.WebGLContext) *common.Shader {
	// Shader for three axes - X(RED) & Y(GREEN) - for visual reference,
	//   with auto-binded (Proj * View * Model) matrix and XY coordinates
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		varying   vec2 v_xy;		// (varying) XY coordinates
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_xy = xy;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec2 v_xy;			// (varying) XY coordinates
		void main() {
			if (v_xy.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else               gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")     // automatic binding of Proj*View*Model matrix
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords") // automatic binding of point coordinates
	shader.SetThingsToDraw("LINES")                                 // only for axis LINES
	return shader
}

func NewShader_Basic(wctx *common.WebGLContext) *common.Shader {
	// Shader with auto-binded color and (Proj * View * Model) matrix
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// color
		void main() { 
			gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")     // automatic binding of Proj*View*Model matrix
	shader.InitBindingForUniform("color", "vec3", "material.color") // automatic binding of material color
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords") // automatic binding of point coordinates
	shader.SetThingsToDraw("LINES", "TRIANGLES")                    // can be used for drawing either
	return shader
}

func NewShader_BasicTexture(wctx *common.WebGLContext) *common.Shader {
	// Shader with auto-binded color and (Proj * View * Model) matrix
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		attribute vec2 uv;			// UV coordinates
		varying vec2 v_uv;			// (varying) UV coordinates
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_uv = uv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_uv;			// (varying) UV coordinates
		void main() { 
			gl_FragColor = texture2D(text, v_uv);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")           // automatic binding of Proj*View*Model matrix
	shader.InitBindingForUniform("text", "sampler2D", "material.texture") // automatic binding of texture sampler (unit:0)
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords")       // automatic binding of point coordinates
	shader.InitBindingForAttribute("uv", "vec2", "geometry.textuv")       // automatic binding of texture UV coordinates
	shader.SetThingsToDraw("LINES", "TRIANGLES")                          // can be used for drawing either
	return shader
}

func NewShader_InstancePoseColor(wctx *common.WebGLContext) *common.Shader {
	// Shader with instance pose, for rendering multiple instances of a same geometry
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		attribute vec2 ixy;			// instance pose : XY translation
		attribute vec3 icolor;		// instance pose : color
		varying   vec3 v_color;		// (varying) color
		void main() {
			vec3 new_pos = pvm * vec3(xy.x + ixy.x, xy.y + ixy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_color = icolor;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying   vec3 v_color;		// (varying) color
		void main() { 
			gl_FragColor = vec4(v_color, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat3", "renderer.pvm")           // automatic binding of Proj*View*Model matrix
	shader.InitBindingForAttribute("xy", "vec2", "geometry.coords")       // automatic binding of point coordinates
	shader.InitBindingForAttribute("ixy", "vec2", "instance.pose:5:0")    // automatic binding of instance pose
	shader.InitBindingForAttribute("icolor", "vec3", "instance.pose:5:2") // automatic binding of instance color
	shader.SetThingsToDraw("LINES", "TRIANGLES")                          // can be used for drawing either
	return shader
}
