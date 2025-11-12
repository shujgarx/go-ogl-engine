package gfx

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	FOV float32
	Aspect float32
	Near float32
	Far float32
}

func (c *Camera) MVP(pos, rot, scale [3]float32) mgl.Mat4 {
	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(pos[0], pos[1], pos[2]))
	// simple Y rotation only, for demo
	model = model.Mul4(mgl.HomogRotate3DY(rot[1]))
	model = model.Mul4(mgl.Scale3D(scale[0], scale[1], scale[2]))

	proj := mgl.Perspective(mgl.DegToRad(c.FOV), c.Aspect, c.Near, c.Far)
	view := mgl.LookAtV(mgl.Vec3{0, 0, 3}, mgl.Vec3{0,0,0}, mgl.Vec3{0,1,0})
	return proj.Mul4(view).Mul4(model)
}

func Deg(x float32) float32 { return float32(float64(x) * math.Pi / 180.0) }
