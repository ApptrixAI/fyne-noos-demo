//go:build tinygo || noos

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/embedded"

	xEmbedded "fyne.io/x/fyne/driver/embedded"
)

func setup(a fyne.App) embedded.Driver {
	d := xEmbedded.NewTinyGoDriver()
	app.SetDriverDetails(a, d)
	return d
}
