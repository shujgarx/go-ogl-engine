package gfx

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
)

type Mesh struct {
	VAO, VBO, EBO uint32
	IndexCount    int32
}

func NewMesh(vertices []float32, indices []uint32) *Mesh {
	m := &Mesh{}
	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.EBO)

	gl.BindVertexArray(m.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// layout: pos(3), normal(3), tex(2)
	stride := int32(8 * 4)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(0)))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(12)))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(24)))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	m.IndexCount = int32(len(indices))
	return m
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}
