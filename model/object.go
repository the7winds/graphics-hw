package model

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Object represents a particular object to render with buffers and stuff
type Object struct {
	m *Model
	M mgl32.Mat4
}

// NewObject create a particular object by a model
func (m *Model) NewObject() *Object {
	object := new(Object)
	object.m = m
	object.M = mgl32.Ident4()

	return object
}

func (obj *Object) Draw(programID uint32) {
	gl.UseProgram(programID)

	gl.BindVertexArray(obj.m.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, obj.m.vertexBuffer)
	Vertex := uint32(gl.GetAttribLocation(programID, gl.Str("Vertex\x00")))
	gl.EnableVertexAttribArray(Vertex)
	gl.VertexAttribPointer(Vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, obj.m.texCoordBuffer)
	if loc := gl.GetAttribLocation(programID, gl.Str("TPos\x00")); loc != -1 {
		TPos := uint32(loc)
		gl.EnableVertexAttribArray(TPos)
		gl.VertexAttribPointer(TPos, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	}

	gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("M\x00")), 1, false, &obj.M[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.m.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, int32(len(obj.m.faces)), gl.UNSIGNED_SHORT, nil)
}

func (obj *Object) Model() *Model {
	return obj.m
}
