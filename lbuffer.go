package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type LBuffer struct {
	programID      uint32
	fbo            uint32
	torches        []*Torch
	lightTexture   uint32
	volumesTexture uint32
	depthBuffer    uint32
}

func (lbuffer *LBuffer) init() error {
	lbuffer.programID = newProgram("shaders/lbuffer/vertex.glsl", "shaders/lbuffer/fragment.glsl")

	lbuffer.loadLight()

	gl.GenFramebuffers(1, &lbuffer.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, lbuffer.fbo)

	gl.GenTextures(1, &lbuffer.lightTexture)
	gl.BindTexture(gl.TEXTURE_2D, lbuffer.lightTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 800, 800, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, lbuffer.lightTexture, 0)

	gl.GenTextures(1, &lbuffer.volumesTexture)
	gl.BindTexture(gl.TEXTURE_2D, lbuffer.volumesTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 800, 800, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT1, lbuffer.volumesTexture, 0)

	gl.GenRenderbuffers(1, &lbuffer.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, lbuffer.depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, 800, 800)
	gl.FramebufferRenderbuffer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, lbuffer.depthBuffer)

	gl.DrawBuffers(3, &([]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.NONE})[0])

	if err := checkGlError("can't init L-Buffer"); err != nil {
		return err
	}

	if errCode := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); errCode != gl.FRAMEBUFFER_COMPLETE {
		errMessage := fmt.Sprintln("can't init L-Buffer (CheckFrameBufferStatus):", errCode)
		return errors.New(errMessage)
	}

	return nil
}

func (lbuffer *LBuffer) free() {
	gl.DeleteFramebuffers(1, &lbuffer.fbo)
}

func (lbuffer *LBuffer) loadLight() {
	sphereModel := NewModel("objects/icosphere.obj")

	torch := NewTorch(sphereModel)
	torch.color = mgl32.Vec4{1, 0, 0, 0}
	torch.Scale(2)
	torch.Move(0, 1, 0)
	lbuffer.torches = append(lbuffer.torches, torch)

	torch = NewTorch(sphereModel)
	torch.color = mgl32.Vec4{0, 1, 0, 0}
	torch.Scale(2)
	torch.Move(0, 1, 1)
	lbuffer.torches = append(lbuffer.torches, torch)

	torch = NewTorch(sphereModel)
	torch.color = mgl32.Vec4{0, 0, 1, 0}
	torch.Scale(2)
	torch.Move(0, 1, -1)
	lbuffer.torches = append(lbuffer.torches, torch)

	n := 100
	for i := 0; i < n; i++ {
		torch = NewTorch(sphereModel)
		torch.color = mgl32.Vec4{rand.Float32(), rand.Float32(), rand.Float32(), 0}
		torch.Scale(2)
		torch.Move(0, 1, 0)
		lbuffer.torches = append(lbuffer.torches, torch)
	}
}

func (lbuffer *LBuffer) render(gbuffer *GBuffer, PV *mgl32.Mat4) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, lbuffer.fbo)
	gl.Viewport(0, 0, 800, 800)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.Disable(gl.DEPTH_TEST)
	defer gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0, 0, 0, 0)

	gl.UseProgram(lbuffer.programID)

	gl.UniformMatrix4fv(gl.GetUniformLocation(lbuffer.programID, gl.Str("PV\x00")), 1, false, &PV[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.colorTexture)
	gl.Uniform1i(gl.GetUniformLocation(lbuffer.programID, gl.Str("TexColor\x00")), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.normaTexture)
	gl.Uniform1i(gl.GetUniformLocation(lbuffer.programID, gl.Str("TexNorma\x00")), 1)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.depthTexture)
	gl.Uniform1i(gl.GetUniformLocation(lbuffer.programID, gl.Str("TexDepth\x00")), 2)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.ONE, gl.ONE)
	defer gl.Disable(gl.BLEND)

	gl.Enable(gl.CULL_FACE)
	defer gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)
	defer gl.CullFace(gl.BACK)

	for _, torch := range lbuffer.torches {
		torch.Animate()
		torch.draw(lbuffer.programID)
	}

	return checkGlError("can't render L-Buffer")
}
