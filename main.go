package main

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const width = 800
const height = 800

type Application struct {
	window *glfw.Window
	scene  Scene
}

func (app *Application) glfwInit() error {
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "deferred", nil, nil)

	if err != nil {
		return err
	}

	app.window = window
	app.window.MakeContextCurrent()

	return nil
}

func (app *Application) glInit() error {
	if err := gl.Init(); err != nil {
		return err
	}

	gl.ClearColor(0.0, 0.1, 0.2, 1.0)

	return nil
}

func (app *Application) scenePrepare() {
	app.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		app.scene.keyCallback(w, key, scancode, action, mods)
	})

	app.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		app.scene.mouseButtonCallback(w, button, action, mod)
	})

	app.window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		app.scene.cursorPosCallback(w, xpos, ypos)
	})

	app.scene.loadModel()
}

func (app *Application) run() {
	for !app.window.ShouldClose() {
		app.scene.render()
		app.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func main() {
	var app Application

	if err := app.glfwInit(); err != nil {
		log.Fatalln(err)
		return
	}

	if err := app.glInit(); err != nil {
		log.Fatalln(err)
		return
	}

	app.scenePrepare()
	app.run()
}
