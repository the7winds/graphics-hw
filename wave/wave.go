package wave

import (
	"errors"
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/the7winds/graphics-hw/model"
	"github.com/the7winds/graphics-hw/program"
	"github.com/the7winds/graphics-hw/utils"
)

const (
	texWidth  = 800
	texHeight = 800
)

type Wave struct {
	fbo uint32
	tex uint32

	programID uint32

	Pause bool

	model     *model.Model
	particles []*particle
	extra     []*particle
}

func (w *Wave) genParticles() {
	delta := 1.57
	n := 20
	for i := 0; i < n; i++ {
		p := 2 * math.Pi * float64(i) / float64(n-1)
		x := 0.9 * math.Sin(1*p+delta)
		y := 0.9 * math.Sin(2*p)
		w.particles = append(w.particles, newRandomStampedParticle(float32(x), float32(y), w.model))
	}
}

func (w *Wave) initBuffers() {
	gl.GenFramebuffers(1, &w.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, w.fbo)

	gl.GenTextures(1, &w.tex)
	gl.BindTexture(gl.TEXTURE_2D, w.tex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16F, texWidth, texHeight, 0, gl.RED, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, w.tex, 0)
	gl.DrawBuffers(1, &([]uint32{gl.COLOR_ATTACHMENT0})[0])

	if err := utils.CheckGlError("can't init wave-buffer"); err != nil {
		panic(err)
	}

	if errCode := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); errCode != gl.FRAMEBUFFER_COMPLETE {
		errMessage := fmt.Sprintln("can't init wave-buffer (CheckFrameBufferStatus):", errCode)
		panic(errors.New(errMessage))
	}
}

func New() *Wave {
	w := new(Wave)
	w.model = model.Load("objects/screen.obj")
	w.genParticles()
	w.initBuffers()
	w.programID = program.New("shaders/wave/vertex.glsl", "shaders/wave/fragment.glsl")

	return w
}

func (w *Wave) Free() {
	w.model.Free()
	gl.DeleteTextures(1, &w.tex)
	gl.DeleteFramebuffers(1, &w.fbo)
}

func (w *Wave) Animate() {
	if w.Pause {
		return
	}

	for _, p := range w.particles {
		p.animate()
	}

	tmp := w.extra
	w.extra = nil

	for _, p := range tmp {
		if p.time+1 != p.period {
			w.extra = append(w.extra, p)
		}
	}

	for _, p := range w.extra {
		p.animate()
	}
}

func (w *Wave) Render() error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, w.fbo)
	gl.Viewport(0, 0, texWidth, texHeight)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(w.programID)

	gl.Disable(gl.DEPTH_TEST)
	defer gl.Enable(gl.DEPTH_TEST)

	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)

	for _, p := range w.particles {
		p.draw(w.programID)
	}

	for _, p := range w.extra {
		p.draw(w.programID)
	}

	return utils.CheckGlError("can't render wave-Buffer")
}

func (w *Wave) Tex() uint32 {
	return w.tex
}

func (w *Wave) EmitAt(x, y float32) {
	if len(w.extra) > 0 {
		last := w.extra[len(w.extra)-1]
		if math.Abs(float64(last.x-x)) < 0.02 && math.Abs(float64(last.y-y)) < 0.02 {
			if last.time < last.period/10 {
				return
			}
		}
	}
	w.extra = append(w.extra, newParticle(x, y, w.model))
}

func (w *Wave) AutogenOnly() {
	w.extra = nil
	w.genParticles()
}

func (w *Wave) UserOnly() {
	w.particles = nil
	w.Pause = false
}
