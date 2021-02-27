package webgl3d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_3DAxes(wctx *common.WebGLContext) *common.Shader {
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;			// Projection * View * Model matrix
		attribute vec3 xyz;			// XYZ coordinates
		varying vec3 v_xyz;			// (varying) XYZ coordinates
		void main() {
			gl_Position = pvm * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			v_xyz = xyz;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_xyz;			// (varying) XYZ coordinates
		void main() {
			if      (v_xyz.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else if (v_xyz.y != 0.0) gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
			else                     gl_FragColor = vec4(0.6, 0.6, 1.0, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat4", "renderer.pvm")      // automatic binding of (Proj * View * Models) matrix
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of vertex coordinates
	shader.SetThingsToDraw("LINES")                                  // only for axis LINES
	return shader
}

func NewShader_ColorOnly(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;			// Projection * View * Model matrix
		attribute vec3 xyz;			// XYZ coordinates
		void main() {
			gl_Position = pvm * vec4(xyz, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// single color
		void main() { 
			gl_FragColor = vec4(color.rgb, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat4", "renderer.pvm")      // automatic binding of (Proj * View * Models) matrix
	shader.InitBindingForUniform("color", "vec3", "material.color")  // automatic binding of material color
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of point XYZ coordinates
	shader.SetThingsToDraw("LINES", "TRIANGLES")                     // can be used for drawing either
	return shader
}

func NewShader_NormalColor(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		uniform mat3 light;			// directional light ([0]:direction, [1]:color, [2]:ambient) COLUMN-MAJOR!
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec3 nor;			// normal vector
		varying vec3 v_light;   	// (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// material color
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			gl_FragColor = vec4(color.rgb * v_light, color.a);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")    // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vwmd", "mat4", "renderer.vwmd")    // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("color", "vec4", "material.color")  // automatic binding of material color
	shader.InitBindingForUniform("light", "mat3", "lighting.dlight") // automatic binding of directional lighting
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("nor", "vec3", "geometry.normal") // automatic binding of point normal vectors
	shader.SetThingsToDraw("TRIANGLES")                              // can be used for drawing TRIANGLES
	return shader
}

func NewShader_TextureOnly(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + UV + NORMAL) Geometry & (TEXTURE) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec2 tuv;			// texture coordinates
		varying vec2 v_tuv;			// (varying) texture coordinates
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			v_tuv = tuv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_tuv;			// (varying) texture coordinates
		void main() { 
			gl_FragColor = texture2D(text, v_tuv);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")         // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vwmd", "mat4", "renderer.vwmd")         // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("text", "sampler2D", "material.texture") // automatic binding of texture sampler (unit:0)
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords")      // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("tuv", "vec2", "geometry.textuv")      // automatic binding of point UV coordinates (texture)
	shader.SetThingsToDraw("TRIANGLES")                                   // can be used for drawing TRIANGLES
	return shader
}

func NewShader_NormalTexture(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + UV + NORMAL) Geometry & (TEXTURE) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		uniform mat3 light;			// directional light ([0]:direction, [1]:color, [2]:ambient) COLUMN-MAJOR!
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec2 tuv;			// texture coordinates
		attribute vec3 nor;			// normal vector
		varying vec2 v_tuv;			// (varying) texture coordinates
		varying vec3 v_light;		// (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			v_tuv = tuv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_tuv;			// (varying) texture coordinates
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			vec4 color = texture2D(text, v_tuv);
			gl_FragColor = vec4(color.rgb * v_light, color.a);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")         // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vwmd", "mat4", "renderer.vwmd")         // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("light", "mat3", "lighting.dlight")      // automatic binding of directional lighting
	shader.InitBindingForUniform("text", "sampler2D", "material.texture") // automatic binding of texture sampler (unit:0)
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords")      // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("tuv", "vec2", "geometry.textuv")      // automatic binding of point UV coordinates (texture)
	shader.InitBindingForAttribute("nor", "vec3", "geometry.normal")      // automatic binding of point normal vector
	shader.SetThingsToDraw("TRIANGLES")                                   // can be used for drawing TRIANGLES
	return shader
}

func NewShader_InstancePoseColor(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec3 nor;			// normal vector
		attribute vec3 ixyz;		// instance pose : XYZ translation
		attribute vec3 icolor;		// instance pose : color
		uniform mat3 light;			// [0]: direction, [1]: color, [2]: ambient_color   (column-major)
		varying vec3 v_color;    	// (varying) instance color
		varying vec3 v_light;    	// (varying) lighting intensity
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x + ixyz[0], xyz.y + ixyz[1], xyz.z + ixyz[2], 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			v_color = icolor;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_color;		// (varying) instance color
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			gl_FragColor = vec4(v_color * v_light, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")         // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vwmd", "mat4", "renderer.vwmd")         // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("light", "mat3", "lighting.dlight")      // automatic binding of directional lighting
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords")      // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("nor", "vec3", "geometry.normal")      // automatic binding of point normal vectors
	shader.InitBindingForAttribute("ixyz", "vec3", "instance.pose:6:0")   // automatic binding of instance position
	shader.InitBindingForAttribute("icolor", "vec3", "instance.pose:6:3") // automatic binding of instance color
	shader.SetThingsToDraw("TRIANGLES")                                   // can be used for drawing TRIANGLES
	return shader
}
