package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"git.ddd.rip/ptrcnull/watchface/framebuffer"
)

var Gray = color.RGBA{R: 60, G: 60, B: 60, A: 60}
var LightGray = color.RGBA{R: 150, G: 150, B: 150, A: 150}
var White = color.RGBA{R: 255, G: 255, B: 255, A: 255}

func addLabel(img draw.Image, face font.Face, rect image.Rectangle, label string) {
	point := fixed.Point26_6{X: fixed.Int26_6(rect.Min.X * 64), Y: fixed.Int26_6(rect.Max.Y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(LightGray),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}

var SecondClock = image.Rect(62, 162, 298, 198)
var MinuteClock = image.Rect(108, 162, 253, 198)
var Battery = image.Rect(90, 240, 150, 240+32)

var notoSans *truetype.Font

func sized(size float64) font.Face {
	return truetype.NewFace(notoSans, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
		DPI:     0,
	})
}

func main() {
	fontData, _ := ioutil.ReadFile("/usr/share/fonts/noto/NotoSansMono-Regular.ttf")
	notoSans, _ = freetype.ParseFont(fontData)

	fb, _ := framebuffer.Open("/dev/fb0")
	Fill(fb, fb.Bounds(), color.RGBA{})

	face := Face{
		//tmp: &framebuffer.SimpleRGBA{
		//	Pixels: make([]uint8, fb.Xres*fb.Yres*4),
		//	Stride: fb.Yres * 4,
		//	Xres:   fb.Xres,
		//	Yres:   fb.Yres,
		//},
		tmp: image.NewRGBA(fb.Bounds()),
		fb:  fb.SimpleRGBA,
	}

	simple := true
	go face.MinuteClock(time.Now())
	go face.Battery()

	ticker := time.NewTicker(time.Second)
	for {
		t := time.Now()
		if simple {
			if t.Second() == 0 {
				go face.MinuteClock(t)
				go face.Battery()
			}
		} else {
			go face.SecondClock(t)
			if t.Second() == 0 {
				go face.Battery()
			}
		}
		<-ticker.C
	}
}

type Face struct {
	tmp draw.Image
	fb  draw.Image
}

func (f *Face) Battery() {
	Fill(f.tmp, Battery, color.RGBA{})
	addLabel(f.tmp, sized(32), Battery, getBattery())
	Copy(f.fb, f.tmp, Battery)
}

func (f *Face) SecondClock(t time.Time) {
	Fill(f.tmp, SecondClock, color.RGBA{})
	addLabel(f.tmp, sized(48), SecondClock, t.Format("15:04:05"))
	Copy(f.fb, f.tmp, SecondClock)
}

func (f *Face) MinuteClock(t time.Time) {
	Fill(f.tmp, MinuteClock, color.RGBA{})
	addLabel(f.tmp, sized(48), MinuteClock, t.Format("15:04"))
	Copy(f.fb, f.tmp, MinuteClock)
}

//func Loop(duration time.Duration, handler func()) {
//	for {
//		go handler()
//		time.Sleep(duration)
//	}
//}

func getBattery() string {
	res, _ := ioutil.ReadFile("/sys/class/power_supply/battery/capacity")
	value := strings.Trim(string(res), "\n")
	if len(value) == 1 {
		value = "0" + value
	}
	if value == "100" {
		return "uwu"
	}
	return value + "%"
}
