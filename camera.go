package main

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

	P mgl32.Mat4

	PV mgl32.Mat4
}

func (camera *Camera) init(eye mgl32.Vec3) {
	camera.eye = eye
	camera.verticalAngle = 0
	camera.horizontalAngle = 1
	camera.dHorizontal = 0.01
	camera.dVertical = 0.01

	camera.P = mgl32.Perspective(1.4/2, 4.0/3.0, 0.1, 100)

	camera.update()
}

func (camera *Camera) update() {
	ha := float64(camera.horizontalAngle)
	va := float64(camera.verticalAngle)

	x := float32(math.Sin(ha) * math.Cos(va))
	y := float32(math.Sin(ha) * math.Sin(va))
	z := float32(math.Cos(ha))

	camera.dir = mgl32.Vec3{x, y, z}.Normalize()
	camera.ortho = mgl32.Vec3.Cross(camera.dir, mgl32.Vec3{0, 1, 0}).Normalize()

	V := mgl32.LookAtV(camera.eye, mgl32.Vec3.Add(camera.eye, camera.dir), mgl32.Vec3{0, 1, 0})

	camera.PV = mgl32.Mat4.Mul4(camera.P, V)
}

func (camera *Camera) configure(eye, dir mgl32.Vec3) {
	camera.eye = eye
	camera.dir = dir
	camera.update()
}

func (camera *Camera) rotateDeltaUpVertical() {
	camera.verticalAngle += camera.dVertical
	camera.update()
}

func (camera *Camera) moveEyeForward() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.dir.Mul(0.01))
	camera.update()
}

func (camera *Camera) moveEyeBackward() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.dir.Mul(-0.01))
	camera.update()
}

func (camera *Camera) moveEyeRight() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.ortho.Mul(0.01))
	camera.update()
}

func (camera *Camera) moveEyeLeft() {
	camera.eye = mgl32.Vec3.Add(camera.eye, camera.ortho.Mul(-0.01))
	camera.update()
}
