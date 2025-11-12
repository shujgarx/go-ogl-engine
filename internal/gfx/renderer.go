package gfx

import (
	"io/ioutil"
	"path/filepath"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	Program uint32
	Cam     Camera
	UniformMVP int32
}

func NewRenderer(assetDir string, width, height int) (*Renderer, error) {
	vsSrc, _ := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "basic.vert"))
	fsSrc, _ := ioutil.ReadFile(filepath.Join(assetDir, "shaders", "basic.frag"))

	vs, err := CompileShader(string(vsSrc), gl.VERTEX_SHADER)
	if err != nil { return nil, err }
	fs, err := CompileShader(string(fsSrc), gl.FRAGMENT_SHADER)
	if err != nil { return nil, err }
	prog, err := LinkProgram(vs, fs)
	if err != nil { return nil, err }

	r := &Renderer{Program: prog}
	gl.UseProgram(r.Program)
	r.UniformMVP = gl.GetUniformLocation(r.Program, gl.Str("uMVP\x00"))
	gl.UseProgram(0)

	r.Cam = Camera{FOV: 60, Aspect: float32(width)/float32(height), Near: 0.1, Far: 100.0}
	return r, nil
}
