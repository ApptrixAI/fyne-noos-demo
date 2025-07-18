//go:build (tinygo || noos) && nodisplay

package main

import (
	"image"
)

func nextKey() uint16 {
	return noKey
}

func refresh(_ image.Image) {
}

func screenSize() (uint64, uint64) {
	return 320, 240
}
