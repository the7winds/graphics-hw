package wave

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/the7winds/graphics-hw/model"
)

type particle struct {
	x, y   float32
	stamp  int64
	time   int
	period int
	obj    *model.Object
}

func newRandomStampedParticle(x, y float32, m *model.Model) *particle {
	p := new(particle)
	p.obj = m.NewObject()
	p.obj.M = mgl32.Translate3D(x, y, 0).Mul4(p.obj.M)
	p.period = 120
	p.time = rand.Int() % 120
	return p
}

func newParticle(x, y float32, m *model.Model) *particle {
	p := new(particle)
	p.obj = m.NewObject()
	p.x, p.y = x, y
	p.obj.M = mgl32.Translate3D(x, y, 0).Mul4(p.obj.M)
	p.period = 120
	p.time = 0
	return p
}

func (p *particle) animate() {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if now-p.stamp > 20 {
		p.stamp = now
		p.time = (p.time + 1) % p.period
	}
}

func (p *particle) draw(programID uint32) {
	gl.Uniform1f(gl.GetUniformLocation(programID, gl.Str("Time\x00")), float32(p.time)/float32(p.period))
	p.obj.Draw(programID)
}
