package main

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GBuffer struct {
	programID    uint32
	fbo          uint32
	objects      []*Object
	colorTexture uint32
	depthTexture uint32
	normaTexture uint32
}

func (gbuffer *GBuffer) init() error {
	gbuffer.programID = newProgram("shaders/gbuffer/vertex.glsl", "shaders/gbuffer/fragment.glsl")
	gbuffer.loadScene()

	gl.GenFramebuffers(1, &gbuffer.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, gbuffer.fbo)

	gl.GenTextures(1, &gbuffer.colorTexture)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.colorTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 800, 800, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gbuffer.colorTexture, 0)

	gl.GenTextures(1, &gbuffer.normaTexture)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.normaTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 800, 800, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gbuffer.normaTexture, 0)

	gl.GenTextures(1, &gbuffer.depthTexture)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.depthTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, 800, 800, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gbuffer.depthTexture, 0)

	gl.DrawBuffers(3, &([]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.NONE})[0])

	if err := checkGlError("can't init G-Buffer"); err != nil {
		return err
	}

	if errCode := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); errCode != gl.FRAMEBUFFER_COMPLETE {
		errMessage := fmt.Sprintln("can't init G-Buffer (CheckFrameBufferStatus):", errCode)
		return errors.New(errMessage)
	}

	return nil
}

func (gbuffer *GBuffer) loadScene() {
	plane := NewModel("objects/sponza.obj").NewObject()
	plane.color = mgl32.Vec4{1, 1, 1, 1}

	sphere := NewModel("objects/stanford_bunny.obj").NewObject()
	sphere.color = mgl32.Vec4{1, 0, 1, 1}
	sphere.M = sphere.M.Mul4(mgl32.Scale3D(3, 3, 3))

	gbuffer.objects = append(gbuffer.objects, plane, sphere)
}

func (gbuffer GBuffer) render(PV *mgl32.Mat4) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, gbuffer.fbo)
	gl.Viewport(0, 0, 800, 800)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	gl.UseProgram(gbuffer.programID)

	gl.UniformMatrix4fv(gl.GetUniformLocation(gbuffer.programID, gl.Str("PV\x00")), 1, false, &PV[0])

	for _, object := range gbuffer.objects {
		object.draw(gbuffer.programID)
	}

	return checkGlError("can't render G-Buffer")
}
