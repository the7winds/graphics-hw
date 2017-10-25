package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Torch struct {
	model  *Model
	M      mgl32.Mat4
	center mgl32.Vec3
	color  mgl32.Vec4
}

func (model *Model) NewTorch() *Torch {
	torch := new(Torch)
	torch.model = model
	torch.M = mgl32.Ident4()

	return torch
}

func (torch *Torch) draw(programID uint32) {
	gl.UseProgram(programID)

	gl.BindVertexArray(torch.model.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, torch.model.vertexBuffer)
	Vertex := uint32(gl.GetAttribLocation(programID, gl.Str("Vertex\x00")))
	gl.EnableVertexAttribArray(Vertex)
	gl.VertexAttribPointer(Vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// set color
	gl.Uniform4fv(gl.GetUniformLocation(programID, gl.Str("Color\x00")), 1, &torch.color[0])
	// set center
	gl.Uniform3fv(gl.GetUniformLocation(programID, gl.Str("Center\x00")), 1, &torch.center[0])
	// set matrix
	gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("M\x00")), 1, false, &torch.M[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, torch.model.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, int32(len(torch.model.faces)), gl.UNSIGNED_SHORT, nil)
}
