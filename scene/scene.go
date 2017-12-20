package scene

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/the7winds/graphics-hw/consts"
	"github.com/the7winds/graphics-hw/model"
	"github.com/the7winds/graphics-hw/program"
	"github.com/the7winds/graphics-hw/utils"
	"github.com/the7winds/graphics-hw/wave"
)

type Scene struct {
	programID uint32
	fbo       uint32

	// scene object
	surface        *model.Object
	dynamicTexture *wave.Wave
	torch          *mgl32.Vec3

	// cubemap
	cubeMap uint32

	// output
	texRender uint32
	texColor  uint32
	texDepth  uint32
	texNorma  uint32
}

func (scene *Scene) initBuffers() {
	gl.GenFramebuffers(1, &scene.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, scene.fbo)

	gl.GenTextures(1, &scene.texRender)
	gl.BindTexture(gl.TEXTURE_2D, scene.texRender)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, consts.WIDTH, consts.HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, scene.texRender, 0)

	gl.GenTextures(1, &scene.texColor)
	gl.BindTexture(gl.TEXTURE_2D, scene.texColor)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, consts.WIDTH, consts.HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT1, scene.texColor, 0)

	gl.GenTextures(1, &scene.texNorma)
	gl.BindTexture(gl.TEXTURE_2D, scene.texNorma)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, consts.WIDTH, consts.HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT2, scene.texNorma, 0)

	gl.GenTextures(1, &scene.texDepth)
	gl.BindTexture(gl.TEXTURE_2D, scene.texDepth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, consts.WIDTH, consts.HEIGHT, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, scene.texDepth, 0)

	gl.DrawBuffers(4, &([]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2, gl.NONE})[0])

	if err := utils.CheckGlError("can't init scene-buffer"); err != nil {
		panic(err)
	}

	if errCode := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); errCode != gl.FRAMEBUFFER_COMPLETE {
		errMessage := fmt.Sprintln("can't init scene-buffer (CheckFrameBufferStatus):", errCode)
		panic(errors.New(errMessage))
	}
}

func loadTexture(filename string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	image, err := jpeg.Decode(file)

	if err != nil {
		panic(err)
	}

	w := image.Bounds().Size().X
	h := image.Bounds().Size().Y

	bytes := make([]byte, 2*4*w*h)

	p := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			buf := []uint32{r, g, b, a}
			for _, c := range buf {
				binary.LittleEndian.PutUint16(bytes[p:], uint16(c))
				p += 2
			}
		}
	}

	return bytes
}

func (scene *Scene) loadCubeMap() {
	gl.GenTextures(1, &scene.cubeMap)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, scene.cubeMap)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	texNames := []string{
		"environment/cstormydays_bk.jpg",
		"environment/cstormydays_ft.jpg",
		"environment/cstormydays_up.jpg",
		"environment/cstormydays_dn.jpg",
		"environment/cstormydays_lf.jpg",
		"environment/cstormydays_rt.jpg",
	}

	targets := []uint32{
		gl.TEXTURE_CUBE_MAP_POSITIVE_X,
		gl.TEXTURE_CUBE_MAP_NEGATIVE_X,
		gl.TEXTURE_CUBE_MAP_POSITIVE_Y,
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Y,
		gl.TEXTURE_CUBE_MAP_POSITIVE_Z,
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Z,
	}

	for i, filename := range texNames {
		bytes := loadTexture(filename)
		gl.TexImage2D(targets[i],
			0,
			gl.RGBA,
			128,
			128,
			0,
			gl.RGBA,
			gl.UNSIGNED_SHORT,
			gl.Ptr(bytes))
	}

	if err := utils.CheckGlError("can't load cube map"); err != nil {
		panic(err)
	}
}

func New() *Scene {
	scene := new(Scene)
	scene.programID = program.New("shaders/surface/vertex.glsl", "shaders/surface/fragment.glsl")
	scene.dynamicTexture = wave.New()
	scene.surface = model.Load("objects/plane.obj").NewObject()
	scene.initBuffers()
	scene.loadCubeMap()
	return scene
}

func (scene *Scene) Free() {
	gl.DeleteFramebuffers(1, &scene.fbo)
	gl.DeleteTextures(1, &scene.cubeMap)
	gl.DeleteTextures(1, &scene.texRender)
	gl.DeleteTextures(1, &scene.texColor)
	gl.DeleteTextures(1, &scene.texNorma)
	gl.DeleteTextures(1, &scene.texDepth)
	scene.surface.Model().Free()
	scene.dynamicTexture.Free()
}

func (scene *Scene) Render(eye mgl32.Vec3, PV *mgl32.Mat4) error {
	if err := scene.dynamicTexture.Render(); err != nil {
		return err
	}

	if err := scene.modelRender(eye, PV); err != nil {
		return err
	}

	scene.animate()

	return nil
}

func (scene *Scene) animate() {
	scene.dynamicTexture.Animate()
}

func (scene *Scene) modelRender(eye mgl32.Vec3, PV *mgl32.Mat4) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, scene.fbo)
	gl.Viewport(0, 0, consts.WIDTH, consts.HEIGHT)
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(scene.programID)

	gl.UniformMatrix4fv(gl.GetUniformLocation(scene.programID, gl.Str("PV\x00")), 1, false, &PV[0])

	gl.Uniform3fv(gl.GetUniformLocation(scene.programID, gl.Str("CameraPos\x00")), 1, &eye[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, scene.dynamicTexture.Tex())
	gl.Uniform1i(gl.GetUniformLocation(scene.programID, gl.Str("TexHeight\x00")), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, scene.cubeMap)
	gl.Uniform1i(gl.GetUniformLocation(scene.programID, gl.Str("TexEnv\x00")), 1)

	scene.surface.Draw(scene.programID)

	return utils.CheckGlError("can't render scene-buffer")
}

func (scene *Scene) TexRender() uint32 {
	return scene.texRender
}

func (scene *Scene) TexColor() uint32 {
	return scene.texColor
}

func (scene *Scene) TexNorma() uint32 {
	return scene.texNorma
}

func (scene *Scene) TexHeight() uint32 {
	return scene.dynamicTexture.Tex()
}
