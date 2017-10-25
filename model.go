package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Model represents a 3D model in obj format
type Model struct {
	vao          uint32
	vertexBuffer uint32
	indexBuffer  uint32

	vertices []float32
	faces    []uint16
}

func (model *Model) parseModel(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
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

	return nil
}

func (model *Model) configModelBuffers() error {
	gl.GenVertexArrays(1, &model.vao)
	gl.BindVertexArray(model.vao)

	gl.GenBuffers(1, &model.vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, model.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(model.vertices), gl.Ptr(model.vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &model.indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, model.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(model.faces), gl.Ptr(model.faces), gl.STATIC_DRAW)

	if errCode := gl.GetError(); errCode != 0 {
		return errors.New(fmt.Sprintln("gl error: ", errCode))
	}

	return nil
}

// NewModel reads from an obj-file a 3D model
func NewModel(filename string) *Model {
	model := new(Model)

	if err := model.parseModel(filename); err != nil {
		panic(err)
	}

	if err := model.configModelBuffers(); err != nil {
		panic(err)
	}

	return model
}
