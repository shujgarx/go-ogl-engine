package gfx

import (
	"io/ioutil"
	"path/filepath"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	Program uint32
	Cam     Camera

	LightPosition mgl.Vec3
	LightColor    mgl.Vec3

	UniformModel      int32
	UniformView       int32
	UniformProjection int32
	UniformNormal     int32
	UniformLightPos   int32
	UniformLightColor int32
	UniformViewPos    int32
	UniformTexture    int32
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

	r := &Renderer{Program: prog}
	gl.UseProgram(r.Program)
	r.UniformModel = gl.GetUniformLocation(r.Program, gl.Str("uModel\x00"))
	r.UniformView = gl.GetUniformLocation(r.Program, gl.Str("uView\x00"))
	r.UniformProjection = gl.GetUniformLocation(r.Program, gl.Str("uProjection\x00"))
	r.UniformNormal = gl.GetUniformLocation(r.Program, gl.Str("uNormalMatrix\x00"))
	r.UniformLightPos = gl.GetUniformLocation(r.Program, gl.Str("uLightPos\x00"))
	r.UniformLightColor = gl.GetUniformLocation(r.Program, gl.Str("uLightColor\x00"))
	r.UniformViewPos = gl.GetUniformLocation(r.Program, gl.Str("uViewPos\x00"))
	r.UniformTexture = gl.GetUniformLocation(r.Program, gl.Str("uTex\x00"))
	gl.Uniform1i(r.UniformTexture, 0)
	gl.UseProgram(0)

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
	r.LightColor = mgl.Vec3{1.0, 1.0, 1.0}
	return r, nil
}
