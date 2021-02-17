package webgl3d

import (
	"github.com/go4orward/gowebgl/common"
)

func NewShader_3DAxes(wctx *common.WebGLContext) *common.Shader {
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;
		attribute vec3 xyz;
		varying vec3 vxyz;
		void main() {
			gl_Position = pvm * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			vxyz = xyz;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 vxyz;
		void main() {
			if      (vxyz.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else if (vxyz.y != 0.0) gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
			else                    gl_FragColor = vec4(0.6, 0.6, 1.0, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat4", "renderer.pvm")      // automatic binding of (Proj * View * Models) matrix
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of vertex coordinates
	shader.SetThingsToDraw("LINES")                                  // only for axis LINES
	return shader
}

func NewShader_NoLight(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;
		attribute vec3 xyz;
		void main() {
			gl_Position = pvm * vec4(xyz, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;
		void main() { 
			gl_FragColor = vec4(color.rgb, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat4", "renderer.pvm")      // automatic binding of (Proj * View * Models) matrix
	shader.InitBindingForUniform("color", "vec3", "material.color")  // automatic binding of material color
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of point XYZ coordinates
	shader.SetThingsToDraw("LINES", "TRIANGLES")                     // to be used for drawing either
	return shader
}

func NewShader_Basic(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;
		attribute vec3 xyz;
		void main() {
			gl_Position = pvm * vec4(xyz.xyz, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;
		void main() { 
			gl_FragColor = vec4(color.rgb, color.a);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("pvm", "mat4", "renderer.pvm")      // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("color", "vec4", "material.color")  // automatic binding of material color
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of point XYZ coordinates
	shader.SetThingsToDraw("LINES", "TRIANGLES")                     // to be used for drawing either
	return shader
}

func NewShader_BasicNormal(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;
		uniform mat4 vm;
		attribute vec3 xyz;
		attribute vec3 nor;
		uniform mat3 light;		// [0]: direction, [1]: color, [2]: ambient_color   (column-major)
		varying vec3 vlight;    // (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vm * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			float s = sqrt( vm[0][0]*vm[0][0] + vm[0][1]*vm[0][1] + vm[0][2]*vm[0][2]);  // scaling 
			mat3  mvRot = mat3( vm[0][0]/s, vm[0][1]/s, vm[0][2]/s, vm[1][0]/s, vm[1][1]/s, vm[1][2]/s, vm[2][0]/s, vm[2][1]/s, vm[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			vlight = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// material color
		varying vec3 vlight;		// lighting intensity
		void main() { 
			gl_FragColor = vec4(color.rgb * vlight, color.a);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")    // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vm", "mat4", "renderer.vmod")      // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("color", "vec4", "material.color")  // automatic binding of material color
	shader.InitBindingForUniform("light", "mat3", "lighting.dlight") // automatic binding of directional lighting
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords") // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("nor", "vec3", "geometry.normal") // automatic binding of point normal vectors
	shader.SetThingsToDraw("TRIANGLES")                              // to be used for drawing TRIANGLES
	return shader
}

func NewShader_InstancePoseColor(wctx *common.WebGLContext) *common.Shader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;
		uniform mat4 vm;
		attribute vec3 xyz;
		attribute vec3 nor;
		attribute vec3 txyz;	// instance position
		attribute vec3 color;	// instance color
		uniform mat3 light;		// [0]: direction, [1]: color, [2]: ambient_color   (column-major)
		varying vec3 vcolor;    // (varying) instance color
		varying vec3 vlight;    // (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vm * vec4(xyz.x + txyz[0], xyz.y + txyz[1], xyz.z + txyz[2], 1.0);
			float s = sqrt( vm[0][0]*vm[0][0] + vm[0][1]*vm[0][1] + vm[0][2]*vm[0][2]);  // scaling 
			mat3  mvRot = mat3( vm[0][0]/s, vm[0][1]/s, vm[0][2]/s, vm[1][0]/s, vm[1][1]/s, vm[1][2]/s, vm[2][0]/s, vm[2][1]/s, vm[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			vlight = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			vcolor = color;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 vcolor;		// instance color
		varying vec3 vlight;		// lighting intensity
		void main() { 
			gl_FragColor = vec4(vcolor.rgb * vlight, 1.0);
		}`
	shader, _ := common.NewShader(wctx, vertex_shader_code, fragment_shader_code)
	shader.InitBindingForUniform("proj", "mat4", "renderer.proj")        // automatic binding of (Projection) matrix
	shader.InitBindingForUniform("vm", "mat4", "renderer.vmod")          // automatic binding of (View * Models) matrix
	shader.InitBindingForUniform("light", "mat3", "lighting.dlight")     // automatic binding of directional lighting
	shader.InitBindingForAttribute("xyz", "vec3", "geometry.coords")     // automatic binding of point XYZ coordinates
	shader.InitBindingForAttribute("nor", "vec3", "geometry.normal")     // automatic binding of point normal vectors
	shader.InitBindingForAttribute("txyz", "vec3", "instance.pose:6:0")  // automatic binding of instance position
	shader.InitBindingForAttribute("color", "vec3", "instance.pose:6:3") // automatic binding of instance color
	shader.SetThingsToDraw("TRIANGLES")                                  // to be used for drawing TRIANGLES
	return shader
}
