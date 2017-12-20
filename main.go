package main

import (
	"github.com/the7winds/graphics-hw/application"
)

func main() {
	app := application.New()
	app.Run()
	app.Free()
}
