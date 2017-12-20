package screen

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/the7winds/graphics-hw/camera"
	"github.com/the7winds/graphics-hw/consts"
	"github.com/the7winds/graphics-hw/model"
	"github.com/the7winds/graphics-hw/program"
	"github.com/the7winds/graphics-hw/scene"
	"github.com/the7winds/graphics-hw/utils"
)

type screenMode int

const (
	sceneMode screenMode = iota
	colorMode
	normaMode
	heightMode
)

type Screen struct {
	camera        *camera.Camera
	isRotatingNow bool
	xpos          float32
	ypos          float32

	// display screen
	screen    *model.Object
	mode      screenMode
	displayID uint32

	// scene to render
	scene *scene.Scene
}

func New(scene *scene.Scene) *Screen {
	screen := new(Screen)

	screen.scene = scene
	screen.displayID = program.New("shaders/screen/vertex.glsl", "shaders/screen/fragment.glsl")
	screen.screen = model.Load("objects/screen.obj").NewObject()
	screen.camera = camera.New(mgl32.Vec3{5, 5, 5}, -2.2, -0.05)

	if err := utils.CheckGlError("can't load screen"); err != nil {
		panic(err)
	}

	return screen
}

func (screen *Screen) Free() {
	screen.scene.Free()
	screen.screen.Model().Free()
}

func (screen *Screen) Render() error {
	if err := screen.scene.Render(screen.camera.Eye(), &screen.camera.PV); err != nil {
		return err
	}

	return screen.display()
}

func (screen *Screen) display() error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, consts.WIDTH, consts.HEIGHT)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	gl.UseProgram(screen.displayID)

	gl.Uniform1i(gl.GetUniformLocation(screen.displayID, gl.Str("Mode\x00")), int32(screen.mode))

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, screen.scene.TexRender())
	gl.Uniform1i(gl.GetUniformLocation(screen.displayID, gl.Str("TexRender\x00")), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, screen.scene.TexColor())
	gl.Uniform1i(gl.GetUniformLocation(screen.displayID, gl.Str("TexColor\x00")), 1)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, screen.scene.TexNorma())
	gl.Uniform1i(gl.GetUniformLocation(screen.displayID, gl.Str("TexNorma\x00")), 2)

	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, screen.scene.TexHeight())
	gl.Uniform1i(gl.GetUniformLocation(screen.displayID, gl.Str("TexHeight\x00")), 3)

	screen.screen.Draw(screen.displayID)

	return utils.CheckGlError("can't display")
}
