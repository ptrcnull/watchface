package main

import (
	"image"
	"image/color"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const GlyphRatio = 0.85

var notoSans = loadNoto()

func sized(size float64) font.Face {
	return truetype.NewFace(notoSans, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
		DPI:     100,
	})
}

func loadNoto() *truetype.Font {
	fontData, _ := ioutil.ReadFile("/usr/share/fonts/noto/NotoSansMono-Regular.ttf")
	ttf, _ := freetype.ParseFont(fontData)
	return ttf
}

func Text(x, y int, textLength, fontSize float64, color color.RGBA) *textData {
	width := int(textLength * fontSize * GlyphRatio)
	return &textData{
		rect:     image.Rect(x, y, x+width, y+int(fontSize)),
		fontSize: fontSize,
		color:    color,
	}
}

type textData struct {
	rect     image.Rectangle
	fontSize float64
	color    color.RGBA
}

func (f *Face) Text(textData *textData, text string) {
	Fill(f.tmp, textData.rect, color.RGBA{})

	d := &font.Drawer{
		Dst:  f.tmp,
		Src:  image.NewUniform(textData.color),
		Face: sized(textData.fontSize),
		Dot:  fixed.Point26_6{
			X: fixed.Int26_6(textData.rect.Min.X * 64),
			Y: fixed.Int26_6(textData.rect.Max.Y * 64),
		},
	}
	d.DrawString(text)

	Copy(f.fb, f.tmp, textData.rect)
}
