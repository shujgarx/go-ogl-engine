package gfx

import (
	"io/ioutil"
	"path/filepath"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	Program        uint32
	ShadowProgram  uint32
	Cam            Camera
	LightPosition  mgl.Vec3
	LightDirection mgl.Vec3
	LightColor     mgl.Vec3

	LightSpaceMatrix mgl.Mat4

	UniformModel      int32
	UniformView       int32
	UniformProjection int32
	UniformNormal     int32
	UniformLightPos   int32
	UniformLightColor int32
	UniformLightDir   int32
	UniformViewPos    int32
	UniformTexture    int32
	UniformShadowMap  int32
	UniformLightSpace int32

	ShadowUniformModel      int32
	ShadowUniformLightSpace int32

	DepthMapFBO  uint32
	DepthMap     uint32
	ShadowWidth  int
	ShadowHeight int
}

func NewRenderer(assetDir string, width, height int) (*Renderer, error) {
	vsSrc, err := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "basic.vert"))
	if err != nil {
		return nil, err
	}
	fsSrc, err := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "basic.frag"))
	if err != nil {
		return nil, err
	}

	vs, err := CompileShader(string(vsSrc), gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fs, err := CompileShader(string(fsSrc), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := LinkProgram(vs, fs)
	if err != nil {
		return nil, err
	}

	shadowVsSrc, err := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "shadow_depth.vert"))
	if err != nil {
		return nil, err
	}
	shadowFsSrc, err := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "shadow_depth.frag"))
	if err != nil {
		return nil, err
	}
	shadowVS, err := CompileShader(string(shadowVsSrc), gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	shadowFS, err := CompileShader(string(shadowFsSrc), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	shadowProg, err := LinkProgram(shadowVS, shadowFS)
	if err != nil {
		return nil, err
	}

	r := &Renderer{
		Program:          prog,
		ShadowProgram:    shadowProg,
		ShadowWidth:      2048,
		ShadowHeight:     2048,
		LightSpaceMatrix: mgl.Ident4(),
	}

	gl.UseProgram(r.Program)
	r.UniformModel = gl.GetUniformLocation(r.Program, gl.Str("uModel\x00"))
	r.UniformView = gl.GetUniformLocation(r.Program, gl.Str("uView\x00"))
	r.UniformProjection = gl.GetUniformLocation(r.Program, gl.Str("uProjection\x00"))
	r.UniformNormal = gl.GetUniformLocation(r.Program, gl.Str("uNormalMatrix\x00"))
	r.UniformLightPos = gl.GetUniformLocation(r.Program, gl.Str("uLightPos\x00"))
	r.UniformLightColor = gl.GetUniformLocation(r.Program, gl.Str("uLightColor\x00"))
	r.UniformLightDir = gl.GetUniformLocation(r.Program, gl.Str("uLightDir\x00"))
	r.UniformViewPos = gl.GetUniformLocation(r.Program, gl.Str("uViewPos\x00"))
	r.UniformTexture = gl.GetUniformLocation(r.Program, gl.Str("uTex\x00"))
	r.UniformShadowMap = gl.GetUniformLocation(r.Program, gl.Str("uShadowMap\x00"))
	r.UniformLightSpace = gl.GetUniformLocation(r.Program, gl.Str("uLightSpaceMatrix\x00"))
	gl.Uniform1i(r.UniformTexture, 0)
	gl.Uniform1i(r.UniformShadowMap, 1)
	gl.UseProgram(0)

	gl.UseProgram(r.ShadowProgram)
	r.ShadowUniformModel = gl.GetUniformLocation(r.ShadowProgram, gl.Str("uModel\x00"))
	r.ShadowUniformLightSpace = gl.GetUniformLocation(r.ShadowProgram, gl.Str("uLightSpaceMatrix\x00"))
	gl.UseProgram(0)

	gl.GenFramebuffers(1, &r.DepthMapFBO)
	gl.GenTextures(1, &r.DepthMap)
	gl.BindTexture(gl.TEXTURE_2D, r.DepthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, int32(r.ShadowWidth), int32(r.ShadowHeight), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	border := []float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &border[0])

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.DepthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, r.DepthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	r.Cam = Camera{
		Position: mgl.Vec3{0, 0, 5},
		Target:   mgl.Vec3{0, 0, 0},
		Up:       mgl.Vec3{0, 1, 0},
		FOV:      60,
		Aspect:   float32(width) / float32(height),
		Near:     0.1,
		Far:      100.0,
	}
	r.LightPosition = mgl.Vec3{2.0, 4.0, 2.0}
	r.LightDirection = mgl.Vec3{-0.4, -1.0, -0.3}.Normalize()
	r.LightColor = mgl.Vec3{1.0, 1.0, 1.0}
	return r, nil
}

func (r *Renderer) Resize(width, height int) {
	if height == 0 {
		height = 1
	}
	r.Cam.Aspect = float32(width) / float32(height)
}

func (r *Renderer) BeginShadowPass() {
	gl.Viewport(0, 0, int32(r.ShadowWidth), int32(r.ShadowHeight))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.DepthMapFBO)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.ShadowProgram)
	gl.UniformMatrix4fv(r.ShadowUniformLightSpace, 1, false, &r.LightSpaceMatrix[0])
}

func (r *Renderer) SubmitShadow(model mgl.Mat4) {
	gl.UniformMatrix4fv(r.ShadowUniformModel, 1, false, &model[0])
}

func (r *Renderer) EndShadowPass() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (r *Renderer) BeginMainPass(width, height int, clear [4]float32) {
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.ClearColor(clear[0], clear[1], clear[2], clear[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(r.Program)
	view := r.Cam.ViewMatrix()
	proj := r.Cam.ProjectionMatrix()
	gl.UniformMatrix4fv(r.UniformView, 1, false, &view[0])
	gl.UniformMatrix4fv(r.UniformProjection, 1, false, &proj[0])
	gl.UniformMatrix4fv(r.UniformLightSpace, 1, false, &r.LightSpaceMatrix[0])

	lp := r.LightPosition
	gl.Uniform3f(r.UniformLightPos, lp[0], lp[1], lp[2])
	lc := r.LightColor
	gl.Uniform3f(r.UniformLightColor, lc[0], lc[1], lc[2])
	ld := r.LightDirection.Normalize()
	gl.Uniform3f(r.UniformLightDir, ld[0], ld[1], ld[2])
	camPos := r.Cam.Position
	gl.Uniform3f(r.UniformViewPos, camPos[0], camPos[1], camPos[2])

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, r.DepthMap)
	gl.ActiveTexture(gl.TEXTURE0)
}

func (r *Renderer) SubmitMesh(model mgl.Mat4, normal mgl.Mat3) {
	gl.UniformMatrix4fv(r.UniformModel, 1, false, &model[0])
	gl.UniformMatrix3fv(r.UniformNormal, 1, false, &normal[0])
}

func (r *Renderer) EndMainPass() {
	gl.UseProgram(0)
}
