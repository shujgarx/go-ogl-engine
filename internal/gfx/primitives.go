package gfx

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func NewCubeMesh() *Mesh {
	v := []float32{
		// positions        // normals         // uvs
		// front
		-0.5, -0.5, 0.5, 0, 0, 1, 0, 0,
		0.5, -0.5, 0.5, 0, 0, 1, 1, 0,
		0.5, 0.5, 0.5, 0, 0, 1, 1, 1,
		-0.5, 0.5, 0.5, 0, 0, 1, 0, 1,
		// back
		-0.5, -0.5, -0.5, 0, 0, -1, 1, 0,
		0.5, -0.5, -0.5, 0, 0, -1, 0, 0,
		0.5, 0.5, -0.5, 0, 0, -1, 0, 1,
		-0.5, 0.5, -0.5, 0, 0, -1, 1, 1,
		// left
		-0.5, -0.5, -0.5, -1, 0, 0, 0, 0,
		-0.5, -0.5, 0.5, -1, 0, 0, 1, 0,
		-0.5, 0.5, 0.5, -1, 0, 0, 1, 1,
		-0.5, 0.5, -0.5, -1, 0, 0, 0, 1,
		// right
		0.5, -0.5, -0.5, 1, 0, 0, 1, 0,
		0.5, -0.5, 0.5, 1, 0, 0, 0, 0,
		0.5, 0.5, 0.5, 1, 0, 0, 0, 1,
		0.5, 0.5, -0.5, 1, 0, 0, 1, 1,
		// top
		-0.5, 0.5, -0.5, 0, 1, 0, 0, 1,
		0.5, 0.5, -0.5, 0, 1, 0, 1, 1,
		0.5, 0.5, 0.5, 0, 1, 0, 1, 0,
		-0.5, 0.5, 0.5, 0, 1, 0, 0, 0,
		// bottom
		-0.5, -0.5, -0.5, 0, -1, 0, 1, 1,
		0.5, -0.5, -0.5, 0, -1, 0, 0, 1,
		0.5, -0.5, 0.5, 0, -1, 0, 0, 0,
		-0.5, -0.5, 0.5, 0, -1, 0, 1, 0,
	}
	i := []uint32{
		0, 1, 2, 0, 2, 3,
		4, 6, 5, 4, 7, 6,
		8, 9, 10, 8, 10, 11,
		12, 14, 13, 12, 15, 14,
		16, 18, 17, 16, 19, 18,
		20, 21, 22, 20, 22, 23,
	}
	return NewMesh(v, i)
}

func NewPlaneMesh(tile float32) *Mesh {
	v := []float32{
		-0.5, 0, 0.5, 0, 1, 0, 0, tile,
		0.5, 0, 0.5, 0, 1, 0, tile, tile,
		0.5, 0, -0.5, 0, 1, 0, tile, 0,
		-0.5, 0, -0.5, 0, 1, 0, 0, 0,
	}
	i := []uint32{0, 1, 2, 2, 3, 0}
	return NewMesh(v, i)
}

func NewSphereMesh(stacks, slices int) *Mesh {
	if stacks < 2 {
		stacks = 2
	}
	if slices < 3 {
		slices = 3
	}
	radius := float32(0.5)
	var vertices []float32
	for stack := 0; stack <= stacks; stack++ {
		phi := float64(stack) / float64(stacks) * math.Pi
		y := float32(math.Cos(phi))
		r := float32(math.Sin(phi))
		for slice := 0; slice <= slices; slice++ {
			theta := float64(slice) / float64(slices) * math.Pi * 2
			x := r * float32(math.Cos(theta))
			z := r * float32(math.Sin(theta))
			u := float32(slice) / float32(slices)
			v := float32(stack) / float32(stacks)
			vertices = append(vertices,
				x*radius,
				y*radius,
				z*radius,
				x,
				y,
				z,
				u,
				v,
			)
		}
	}
	var indices []uint32
	stride := slices + 1
	for stack := 0; stack < stacks; stack++ {
		for slice := 0; slice < slices; slice++ {
			first := uint32(stack*stride + slice)
			second := first + uint32(stride)
			indices = append(indices, first, first+1, second)
			indices = append(indices, second, first+1, second+1)
		}
	}
	return NewMesh(vertices, indices)
}

func NewPyramidMesh() *Mesh {
	apex := mgl.Vec3{0, 0.5, 0}
	base := []mgl.Vec3{
		{-0.5, -0.5, 0.5},
		{0.5, -0.5, 0.5},
		{0.5, -0.5, -0.5},
		{-0.5, -0.5, -0.5},
	}
	var vertices []float32
	var indices []uint32
	index := uint32(0)

	// base
	baseNormal := mgl.Vec3{0, -1, 0}
	baseUV := []mgl.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	for _, tri := range [][]int{{0, 1, 2}, {0, 2, 3}} {
		for _, idx := range tri {
			pos := base[idx]
			uv := baseUV[idx]
			vertices = append(vertices,
				pos[0], pos[1], pos[2],
				baseNormal[0], baseNormal[1], baseNormal[2],
				uv[0], uv[1],
			)
			indices = append(indices, index)
			index++
		}
	}

	// sides
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		edge1 := apex.Sub(base[i])
		edge2 := base[next].Sub(base[i])
		normal := edge2.Cross(edge1).Normalize()
		u1 := float32(i) / 4.0
		u2 := float32(i+1) / 4.0
		uvApex := mgl.Vec2{(u1 + u2) * 0.5, 1}
		uvBase1 := mgl.Vec2{u1, 0}
		uvBase2 := mgl.Vec2{u2, 0}

		verts := []struct {
			pos mgl.Vec3
			uv  mgl.Vec2
		}{
			{apex, uvApex},
			{base[next], uvBase2},
			{base[i], uvBase1},
		}
		for _, v := range verts {
			vertices = append(vertices,
				v.pos[0], v.pos[1], v.pos[2],
				normal[0], normal[1], normal[2],
				v.uv[0], v.uv[1],
			)
			indices = append(indices, index)
			index++
		}
	}
	return NewMesh(vertices, indices)
}

func NewConeMesh(segments int) *Mesh {
	if segments < 3 {
		segments = 3
	}
	height := float32(1.0)
	radius := float32(0.5)
	apex := mgl.Vec3{0, height / 2, 0}
	baseY := -height / 2
	center := mgl.Vec3{0, baseY, 0}
	down := mgl.Vec3{0, -1, 0}

	var vertices []float32
	var indices []uint32
	index := uint32(0)

	for i := 0; i < segments; i++ {
		angle := float64(i) / float64(segments) * math.Pi * 2
		nextAngle := float64(i+1) / float64(segments) * math.Pi * 2
		base1 := mgl.Vec3{radius * float32(math.Cos(angle)), baseY, radius * float32(math.Sin(angle))}
		base2 := mgl.Vec3{radius * float32(math.Cos(nextAngle)), baseY, radius * float32(math.Sin(nextAngle))}

		edge1 := apex.Sub(base1)
		edge2 := base2.Sub(base1)
		normal := edge2.Cross(edge1).Normalize()
		u1 := float32(i) / float32(segments)
		u2 := float32(i+1) / float32(segments)
		uvApex := mgl.Vec2{(u1 + u2) * 0.5, 1}
		uvBase1 := mgl.Vec2{u1, 0}
		uvBase2 := mgl.Vec2{u2, 0}

		verts := []struct {
			pos mgl.Vec3
			uv  mgl.Vec2
		}{
			{apex, uvApex},
			{base2, uvBase2},
			{base1, uvBase1},
		}
		for _, v := range verts {
			vertices = append(vertices,
				v.pos[0], v.pos[1], v.pos[2],
				normal[0], normal[1], normal[2],
				v.uv[0], v.uv[1],
			)
			indices = append(indices, index)
			index++
		}

		// base triangle
		baseVerts := []struct {
			pos mgl.Vec3
			uv  mgl.Vec2
		}{
			{center, mgl.Vec2{(u1 + u2) * 0.5, 1}},
			{base1, mgl.Vec2{u1, 0}},
			{base2, mgl.Vec2{u2, 0}},
		}
		for _, v := range baseVerts {
			vertices = append(vertices,
				v.pos[0], v.pos[1], v.pos[2],
				down[0], down[1], down[2],
				v.uv[0], v.uv[1],
			)
			indices = append(indices, index)
			index++
		}
	}
	return NewMesh(vertices, indices)
}
