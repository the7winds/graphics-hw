package main

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	camera  Camera
	objects []*Object

	isRotatingNow bool
	xpos          float32
	ypos          float32

	programID uint32
}

func (scene *Scene) loadModel() {
	scene.programID = newProgram("shaders/vertex.glsl", "shaders/fragment.glsl")

	sphereModel := NewModel("objects/icosphere.obj")
	sphere := sphereModel.NewObject()
	sphere.color = mgl32.Vec4{1, 1, 1, 1}
	sphere.M = sphere.M.Mul4(mgl32.Translate3D(5, 0, 0))

	sponzaModel := NewModel("objects/sponza.obj")
	sponza := sponzaModel.NewObject()
	sponza.color = mgl32.Vec4{1, 1, 1, 1}

	scene.camera.init(mgl32.Vec3{5, 5, 5})

	scene.objects = append(scene.objects, sponza)

	gl.Enable(gl.DEPTH_TEST)
}

func (scene *Scene) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyUp {
		scene.camera.moveEyeForward()
	} else if key == glfw.KeyDown {
		scene.camera.moveEyeBackward()
	} else if key == glfw.KeyLeft {
		scene.camera.moveEyeLeft()
	} else if key == glfw.KeyRight {
		scene.camera.moveEyeRight()
	}
}

func (scene *Scene) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			scene.isRotatingNow = true
			xpos, ypos := w.GetCursorPos()
			scene.xpos, scene.ypos = float32(xpos), float32(ypos)
		} else {
			scene.isRotatingNow = false
		}
	}
}

func (scene *Scene) cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	if scene.isRotatingNow {
		xpos := float32(xpos)
		ypos := float32(ypos)

		dX := xpos - scene.xpos
		dY := ypos - scene.ypos

		scene.xpos, scene.ypos = xpos, ypos

		scene.camera.rotate(dX, dY)
		scene.camera.update()
	}
}

func (scene *Scene) render() error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(scene.programID)
	gl.UniformMatrix4fv(gl.GetUniformLocation(scene.programID, gl.Str("PV\x00")), 1, false, &scene.camera.PV[0])

	for _, object := range scene.objects {
		object.draw(scene.programID)
	}

	if gl.GetError() != 0 {
		errString := fmt.Sprintln(gl.GetError())
		return errors.New(errString)
	}

	return nil
}
