package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	eye             mgl32.Vec3
	dir             mgl32.Vec3
	ortho           mgl32.Vec3
	verticalAngle   float32
	horizontalAngle float32
	dVertical       float32
	dHorizontal     float32
	dForward        float32

	P mgl32.Mat4

	PV mgl32.Mat4
}

func New(eye mgl32.Vec3, horizontalAngle, verticalAngle float32) *Camera {
	camera := new(Camera)
	camera.eye = eye
	camera.verticalAngle = horizontalAngle
	camera.horizontalAngle = verticalAngle
	camera.dHorizontal = 0.01
	camera.dVertical = 0.01
	camera.dForward = 0.1

	camera.P = mgl32.Perspective(1.4/2, 4.0/3.0, 0.1, 100)

	camera.update()

	return camera
}

func (camera *Camera) update() {
	ha := float64(camera.horizontalAngle)
	va := float64(camera.verticalAngle)

	x := float32(math.Cos(va) * math.Sin(ha))
	y := float32(math.Sin(va))
	z := float32(math.Cos(va) * math.Cos(ha))

	camera.dir = mgl32.Vec3{x, y, z}.Normalize()
	camera.ortho = mgl32.Vec3.Cross(camera.dir, mgl32.Vec3{0, 1, 0}).Normalize()

	V := mgl32.LookAtV(camera.eye, mgl32.Vec3.Add(camera.eye, camera.dir), mgl32.Vec3{0, 1, 0})

	camera.PV = mgl32.Mat4.Mul4(camera.P, V)
}

func (camera *Camera) Rotate(dH, dV float32) {
	camera.verticalAngle += dV * camera.dVertical
	camera.horizontalAngle += dH * camera.dHorizontal
	camera.update()
}

func (camera *Camera) MoveEyeForward() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.dir.Mul(camera.dForward))
	camera.update()
}

func (camera *Camera) MoveEyeBackward() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.dir.Mul(-camera.dForward))
	camera.update()
}

func (camera *Camera) MoveEyeRight() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.ortho.Mul(camera.dForward))
	camera.update()
}

func (camera *Camera) MoveEyeLeft() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.ortho.Mul(-camera.dForward))
	camera.update()
}

func (camera *Camera) Eye() mgl32.Vec3 {
	return camera.eye
}
