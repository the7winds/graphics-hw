package main

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type ScreenMode int

const (
	SCENE ScreenMode = iota
	COLOR
	NORMA
	DEPTH
)

type Scene struct {
	camera        Camera
	isRotatingNow bool
	xpos          float32
	ypos          float32

	// screen
	screen    *Object
	mode      ScreenMode
	displayID uint32

	// l-buffer
	lbuffer LBuffer

	// g-buffer
	gbuffer GBuffer
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
	} else if key == glfw.Key0 {
		scene.mode = SCENE
		fmt.Println("set SCENE mode", scene.mode)
	} else if key == glfw.Key1 {
		scene.mode = COLOR
		fmt.Println("set COLOR mode", scene.mode)
	} else if key == glfw.Key2 {
		scene.mode = NORMA
		fmt.Println("set NORMA mode", scene.mode)
	} else if key == glfw.Key3 {
		scene.mode = DEPTH
		fmt.Println("set DEPTH mode", scene.mode)
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
	}
}

func (scene *Scene) loadModel() error {
	if err := scene.gbuffer.init(); err != nil {
		return err
	}

	if err := scene.lbuffer.init(); err != nil {
		return err
	}

	scene.displayID = newProgram("shaders/screen/vertex.glsl", "shaders/screen/fragment.glsl")
	scene.screen = NewModel("objects/screen.obj").NewObject()

	scene.camera.init(mgl32.Vec3{5, 5, 5}, mgl32.DegToRad(0), mgl32.DegToRad(0))

	return checkGlError("can't load screen")
}

func checkGlError(prefix string) error {
	if errCode := gl.GetError(); errCode != 0 {
		errMessage := fmt.Sprintln(prefix, ":", errCode)
		return errors.New(errMessage)
	}

	return nil
}

func (scene *Scene) render() error {
	if err := scene.gbuffer.render(&scene.camera.PV); err != nil {
		return err
	}

	if err := scene.lbuffer.render(&scene.gbuffer, &scene.camera); err != nil {
		return err
	}

	return scene.display()
}

func (scene *Scene) display() error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, 800, 800)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	gl.UseProgram(scene.displayID)

	gl.UniformMatrix4fv(gl.GetUniformLocation(scene.displayID, gl.Str("PV\x00")), 1, false, &scene.camera.PV[0])
	gl.Uniform1i(gl.GetUniformLocation(scene.displayID, gl.Str("Mode\x00")), int32(scene.mode))

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, scene.gbuffer.colorTexture)
	gl.Uniform1i(gl.GetUniformLocation(scene.displayID, gl.Str("TexColor\x00")), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, scene.gbuffer.normaTexture)
	gl.Uniform1i(gl.GetUniformLocation(scene.displayID, gl.Str("TexNorma\x00")), 1)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, scene.gbuffer.depthTexture)
	gl.Uniform1i(gl.GetUniformLocation(scene.displayID, gl.Str("TexDepth\x00")), 2)

	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, scene.lbuffer.lightTexture)
	gl.Uniform1i(gl.GetUniformLocation(scene.displayID, gl.Str("TexLight\x00")), 3)

	scene.screen.draw(scene.displayID)

	return checkGlError("can't display")
}
