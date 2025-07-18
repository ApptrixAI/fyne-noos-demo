//go:build (tinygo || noos) && !nodisplay

package main

import (
	"image"
	"image/draw"

	"github.com/sago35/tinydisplay"
)

var display *tinydisplay.Client

func nextKey() uint16 {
	return display.GetPressedKey()
}

func refresh(img image.Image) {
	// tinydisplay does not like our *image.NRGBA, so copy it
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.ZP, draw.Src)

	display.SetImage(out)
}

func screenSize() (uint64, uint64) {
	display, _ = tinydisplay.NewClient("127.0.0.1", 9812, 0, 0)
	w, h := display.Size()

	return uint64(w), uint64(h)
}
