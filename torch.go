package main

import (
	"math/rand"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Torch struct {
	model  *Model
	M      mgl32.Mat4
	center mgl32.Vec3
	color  mgl32.Vec4
	R      float32

	prevAnimation uint32
	prevDirection uint32
	direction     mgl32.Vec3
}

func NewTorch(model *Model) *Torch {
	torch := new(Torch)
	torch.model = model
	torch.M = mgl32.Ident4()
	torch.R = 1
	torch.direction = newDirection()

	return torch
}

func newDirection() mgl32.Vec3 {
	gen := func() float32 {
		return (2*rand.Float32() - 1) / 10
	}

	return mgl32.Vec3{gen(), gen(), gen()}
}

func (torch *Torch) draw(programID uint32) {
	gl.UseProgram(programID)

	gl.BindVertexArray(torch.model.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, torch.model.vertexBuffer)

	Vertex := uint32(gl.GetAttribLocation(programID, gl.Str("Vertex\x00")))
	gl.EnableVertexAttribArray(Vertex)
	gl.VertexAttribPointer(Vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// set radius
	gl.Uniform1f(gl.GetUniformLocation(programID, gl.Str("R\x00")), torch.R)
	// set color
	gl.Uniform4fv(gl.GetUniformLocation(programID, gl.Str("Color\x00")), 1, &torch.color[0])
	// set center
	gl.Uniform3fv(gl.GetUniformLocation(programID, gl.Str("Center\x00")), 1, &torch.center[0])
	// set matrix
	gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("M\x00")), 1, false, &torch.M[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, torch.model.indexBuffer)

	gl.DrawElements(gl.TRIANGLES, int32(len(torch.model.faces)), gl.UNSIGNED_SHORT, nil)
}

func (torch *Torch) Scale(s float32) {
	torch.R *= s
	torch.M = mgl32.Scale3D(s, s, s).Mul4(torch.M)
}

func (torch *Torch) Move(dx, dy, dz float32) {
	torch.center = torch.center.Add(mgl32.Vec3{dx, dy, dz})
	torch.M = mgl32.Translate3D(dx, dy, dz).Mul4(torch.M)
}

func (torch *Torch) Animate() {
	now := uint32(time.Now().UnixNano() / int64(time.Millisecond))

	if now-torch.prevAnimation > 20 {
		torch.prevAnimation = now
		torch.Move(torch.direction.X(), torch.direction.Y(), torch.direction.Z())
	}

	if torch.center.X() < -5 || torch.center.X() > 5 ||
		torch.center.Z() < -5 || torch.center.Z() > 5 ||
		torch.center.Y() < -1 || torch.center.Y() > 3 ||
		now-torch.prevDirection > 500 {
		torch.prevDirection = now

		for {
			newDirection := newDirection()
			newCenter := torch.center.Add(newDirection)

			if torch.center.Len() > 20 && torch.center.Len()-newCenter.Len() < 0.01 {
				continue
			}

			torch.direction = newDirection
			torch.Move(torch.direction.X(), torch.direction.Y(), torch.direction.Z())
			return
		}
	}
}
