package screen

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/the7winds/graphics-hw/consts"
)

func (screen *Screen) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyUp:
		screen.camera.MoveEyeForward()
	case glfw.KeyDown:
		screen.camera.MoveEyeBackward()
	case glfw.KeyLeft:
		screen.camera.MoveEyeLeft()
	case glfw.KeyRight:
		screen.camera.MoveEyeRight()
	case glfw.Key0:
		screen.mode = sceneMode
		fmt.Println("set SCENE mode", screen.mode)
	case glfw.Key1:
		screen.mode = colorMode
		fmt.Println("set COLOR mode", screen.mode)
	case glfw.Key2:
		screen.mode = normaMode
		fmt.Println("set NORMA mode", screen.mode)
	case glfw.Key3:
		screen.mode = heightMode
		fmt.Println("set HEIGHT mode", screen.mode)
	}

	if action == glfw.Press {
		switch key {
		case glfw.KeyLeftShift:
			screen.emition = !screen.emition
			if screen.emition {
				screen.scene.SetAutogenMode(false)
				fmt.Println("set USER EMITION mode")
			} else {
				screen.scene.SetAutogenMode(true)
				screen.setPause(false)
				fmt.Println("set AUTO EMITION mode")
			}
		case glfw.KeySpace:
			screen.setPause(!screen.pause)
			if screen.pause {
				fmt.Println("EMITION PAUSED")
			} else {
				fmt.Println("EMITION CONTINUED")
			}
		}
	}
}

func (screen *Screen) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			screen.isRotatingNow = true
			xpos, ypos := w.GetCursorPos()
			screen.xpos, screen.ypos = float32(xpos), float32(ypos)
		} else {
			screen.isRotatingNow = false
		}
	}
}

func (screen *Screen) CursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	if screen.isRotatingNow {
		xpos := float32(xpos)
		ypos := float32(ypos)

		dX := xpos - screen.xpos
		dY := ypos - screen.ypos

		screen.xpos, screen.ypos = xpos, ypos

		screen.camera.Rotate(dX, dY)
	}

	if !screen.pause && screen.emition {
		e := screen.camera.Eye()
		x, y := w.GetCursorPos()
		x = 2*x/consts.WIDTH - 1
		y = 1 - 2*y/consts.HEIGHT
		p := screen.camera.ScreenToWorld(float32(x), float32(y))
		screen.scene.EmitAt(e, p.Sub(e).Normalize())
	}
}
