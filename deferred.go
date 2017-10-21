package main

import (
	"log"
//	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
//	"github.com/go-gl/mathgl/mgl32"
)

const width = 800
const height = 800

func main() {
	/* create window */
	if err := glfw.Init(); err != nil {
		log.Fatalln("can't init glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "deferred", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	/* load GL */
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0.0, 0.0, 0.1, 1.0);

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT);

		window.SwapBuffers();
		glfw.PollEvents();
	}
}

