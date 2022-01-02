package main

import (
	"image"
	"image/color"
	"time"

	"golang.org/x/image/draw"

	"git.ddd.rip/ptrcnull/watchface/framebuffer"
)

var Gray = color.RGBA{R: 60, G: 60, B: 60, A: 60}
var LightGray = color.RGBA{R: 150, G: 150, B: 150, A: 150}
var White = color.RGBA{R: 255, G: 255, B: 255, A: 255}

type App interface {
	Init(*Face)
	Render(*Face)
}

var app App = SimpleClock{}
var launch func(App)

type Face struct {
	tmp draw.Image
	fb  draw.Image
}

func main() {
	fb, _ := framebuffer.Open("/dev/fb0")
	Fill(fb, fb.Bounds(), color.RGBA{})

	face := Face{
		tmp: image.NewRGBA(fb.Bounds()),
		fb:  fb.SimpleRGBA,
	}

	launch = func(newApp App) {
		app = newApp
		app.Init(&face)
	}

	app.Init(&face)

	ticker := time.NewTicker(time.Second)
	for {
		go app.Render(&face)
		<-ticker.C
	}
}
