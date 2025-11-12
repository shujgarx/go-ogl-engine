package gfx

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position mgl.Vec3
	Target   mgl.Vec3
	Up       mgl.Vec3

	FOV    float32
	Aspect float32
	Near   float32
	Far    float32
}

func (c *Camera) ViewMatrix() mgl.Mat4 {
	return mgl.LookAtV(c.Position, c.Target, c.Up)
}

func (c *Camera) ProjectionMatrix() mgl.Mat4 {
	return mgl.Perspective(mgl.DegToRad(c.FOV), c.Aspect, c.Near, c.Far)
}

func (c *Camera) MVP(pos, rot, scale [3]float32) mgl.Mat4 {
	model := ModelMatrix(pos, rot, scale)
	return c.ProjectionMatrix().Mul4(c.ViewMatrix()).Mul4(model)
}

func ModelMatrix(pos, rot, scale [3]float32) mgl.Mat4 {
	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(pos[0], pos[1], pos[2]))
	model = model.Mul4(mgl.HomogRotate3DX(rot[0]))
	model = model.Mul4(mgl.HomogRotate3DY(rot[1]))
	model = model.Mul4(mgl.HomogRotate3DZ(rot[2]))
	model = model.Mul4(mgl.Scale3D(scale[0], scale[1], scale[2]))
	return model
}

func NormalMatrix(model mgl.Mat4) mgl.Mat3 {
	return model.Inv().Transpose().Mat3()
}

func Deg(x float32) float32 { return float32(float64(x) * math.Pi / 180.0) }
