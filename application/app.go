package application

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/the7winds/graphics-hw/consts"
	"github.com/the7winds/graphics-hw/scene"
	"github.com/the7winds/graphics-hw/screen"
	"github.com/the7winds/graphics-hw/utils"
)

type Application struct {
	window *glfw.Window
	screen *screen.Screen
}

func New() *Application {
	app := new(Application)

	if err := app.glfwInit(); err != nil {
		log.Fatalln(err)
	}

	if err := app.glInit(); err != nil {
		log.Fatalln(err)
	}

	if err := utils.CheckGlError("can't init G-Buffer"); err != nil {
		panic(err)
	}

	app.scenePrepare()

	return app
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

	window, err := glfw.CreateWindow(consts.WIDTH, consts.HEIGHT, "waves", nil, nil)

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
	s := scene.New()
	app.screen = screen.New(s)

	app.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		app.screen.KeyCallback(w, key, scancode, action, mods)
	})

	app.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		app.screen.MouseButtonCallback(w, button, action, mod)
	})

	app.window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		app.screen.CursorPosCallback(w, xpos, ypos)
	})
}

func (app *Application) Free() {
	app.screen.Free()
}

func (app *Application) Run() {
	for !app.window.ShouldClose() {
		if err := app.screen.Render(); err != nil {
			panic(err)
		}

		app.window.SwapBuffers()
		glfw.PollEvents()
	}
}
