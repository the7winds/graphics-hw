package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Model represents a 3D model in obj format
type Model struct {
	vertices []float32
	faces    []uint16
}

// Object represents a particular object to render with buffers and stuff
type Object struct {
	model        *Model
	vao          uint32
	vertexBuffer uint32
	indexBuffer  uint32
	M            mgl32.Mat4
	color        mgl32.Vec4
}

// NewModel reads from an obj-file a 3D model
func NewModel(filename string) *Model {
	model := new(Model)
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		var t rune
		fmt.Sscanf(line, "%c", &t)

		if t == 'v' {
			var x, y, z float32
			fmt.Sscanf(line, "v %f %f %f", &x, &y, &z)
			model.vertices = append(model.vertices, x, y, z)
		} else if t == 'f' {
			var a [10]uint16
			n, _ := fmt.Sscanf(line, "f %d %d %d %d %d %d %d %d %d %d", &a[0], &a[1], &a[2], &a[3], &a[4], &a[5], &a[6], &a[7], &a[8], &a[9])
			for i := 1; i < n-1; i++ {
				model.faces = append(model.faces, a[0]-1, a[i]-1, a[i+1]-1)
			}
		}
	}

	return model
}

// NewObject create a particular object by a model
func (model *Model) NewObject() *Object {
	object := new(Object)
	object.model = model
	object.M = mgl32.Ident4()

	gl.GenVertexArrays(1, &object.vao)
	gl.BindVertexArray(object.vao)

	gl.GenBuffers(1, &object.vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, object.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(model.vertices), gl.Ptr(model.vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &object.indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, object.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(model.faces), gl.Ptr(model.faces), gl.STATIC_DRAW)

	return object
}

func (object *Object) draw(programID uint32) {
	gl.UseProgram(programID)

	gl.BindVertexArray(object.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, object.vertexBuffer)
	Vertex := uint32(gl.GetAttribLocation(programID, gl.Str("Vertex\x00")))
	gl.EnableVertexAttribArray(Vertex)
	gl.VertexAttribPointer(Vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// set color
	gl.Uniform4fv(gl.GetUniformLocation(programID, gl.Str("Color\x00")), 1, &object.color[0])

	// set matrix
	gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("M\x00")), 1, false, &object.M[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, object.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, int32(len(object.model.faces)), gl.UNSIGNED_SHORT, nil)
}
