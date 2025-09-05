//go:build !tamago && !tinygo && !noos

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
)

func setup(fyne.App) embedded.Driver {
	return nil
}
