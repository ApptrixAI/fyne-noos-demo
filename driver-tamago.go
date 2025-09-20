//go:build tamago

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	xEmbedded "fyne.io/x/fyne/driver/embedded"
)

func setup(a fyne.App) {
	app.SetDriverDetails(a, xEmbedded.NewUEFIDriver())
}
