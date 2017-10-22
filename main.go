package main

import "log"

import "github.com/go-gl/gl/v4.1-core/gl"
import "github.com/go-gl/glfw/v3.2/glfw"
import "github.com/go-gl/mathgl/mgl32"

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

	gl.ClearColor(0.0, 0.1, 0.2, 1.0)

	// load shader program
	programID := newProgram("shaders/vertex.glsl", "shaders/fragment.glsl")

	model := NewModel("objects/cube.obj")
	object := model.NewObject()
	object.color = mgl32.Vec4{1, 1, 1, 1}

	eye := mgl32.Vec3{5, 5, 5}
	direction := mgl32.Vec3{-1, -1, -1}
	center := eye.Add(direction)

	P := mgl32.Perspective(1.4/2, 4.0/3.0, 0.1, 100)
	V := mgl32.LookAtV(eye, center, mgl32.Vec3{0, 1, 0})
	PV := P.Mul4(V)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyW {
			eye = eye.Add(direction.Mul(0.1))
		} else if key == glfw.KeyS {
			eye = eye.Add(direction.Mul(-0.1))
		} else if key == glfw.KeyUp {
			t := direction.Cross(mgl32.Vec3{0, 1, 0})
			direction = mgl32.HomogRotate3D(mgl32.DegToRad(1), t).Mul4x1(mgl32.Vec4{direction.X(), direction.Y(), direction.Z(), 1}).Vec3()
		} else if key == glfw.KeyDown {
			t := direction.Cross(mgl32.Vec3{0, 1, 0})
			direction = mgl32.HomogRotate3D(mgl32.DegToRad(-1), t).Mul4x1(mgl32.Vec4{direction.X(), direction.Y(), direction.Z(), 1}).Vec3()
		} else if key == glfw.KeyLeft {
			direction = mgl32.Rotate3DY(mgl32.DegToRad(1)).Mul3x1(direction)
		} else if key == glfw.KeyRight {
			direction = mgl32.Rotate3DY(mgl32.DegToRad(-1)).Mul3x1(direction)
		}
		center = eye.Add(direction)
		V = mgl32.LookAtV(eye, center, mgl32.Vec3{0, 1, 0})
		PV = P.Mul4(V)
	})

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(programID)
		gl.UniformMatrix4fv(gl.GetUniformLocation(programID, gl.Str("PV\x00")), 1, false, &PV[0])

		object.draw(programID)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
