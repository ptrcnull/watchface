package main

import (
	"image"
	"image/color"
	"image/draw"
)

func Fill(img draw.Image, rect image.Rectangle, c color.Color) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			img.Set(x, y, c)
		}
	}
}

func Copy(dst, src draw.Image, rect image.Rectangle) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
}
