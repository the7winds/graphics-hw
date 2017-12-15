package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Object represents a particular object to render with buffers and stuff
type Object struct {
	model *Model
	M     mgl32.Mat4
	color mgl32.Vec4
}

// NewObject create a particular object by a model
func (model *Model) NewObject() *Object {
	object := new(Object)
	object.model = model
	object.M = mgl32.Ident4()

	return object
}

func (object *Object) free() {

}

func (object *Object) draw(programID uint32) {
	gl.UseProgram(programID)

	gl.BindVertexArray(object.model.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, object.model.vertexBuffer)
	Vertex := uint32(gl.GetAttribLocation(programID, gl.Str("Vertex\x00")))
	gl.EnableVertexAttribArray(Vertex)
	gl.VertexAttribPointer(Vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	normaPassed := int32(len(object.model.norms))
	gl.Uniform1i(gl.GetUniformLocation(programID, gl.Str("NormaPassed\x00")), normaPassed)

	gl.BindBuffer(gl.ARRAY_BUFFER, object.model.normaBuffer)
	Norma := uint32(gl.GetAttribLocation(programID, gl.Str("Norm\x00")))
	gl.EnableVertexAttribArray(Norma)
	gl.VertexAttribPointer(Norma, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// set color
	gl.Uniform4fv(gl.GetUniformLocation(programID, gl.Str("Color\x00")), 1, &object.color[0])
	// set matrix
	gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("M\x00")), 1, false, &object.M[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, object.model.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, int32(len(object.model.faces)), gl.UNSIGNED_SHORT, nil)
}
