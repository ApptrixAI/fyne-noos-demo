package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed fyne.png
var imgData []byte

func main() {
	a := app.New()
	setup(a)

	w := a.NewWindow("Fyne NOOS Demo")
	setContent(w)
	w.Show()

	a.Run()
}

func setContent(w fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Your name here")

	welcome := widget.NewLabel("Hello NoOS Fyne!\n(press Esc to quit)")
	welcome.Alignment = fyne.TextAlignCenter
	img := canvas.NewImageFromResource(fyne.NewStaticResource("fyne.png", imgData))
	img.SetMinSize(fyne.NewSquareSize(64))

	content := container.NewVBox(welcome, container.NewCenter(img),
		entry,
		widget.NewButton("Greet me", func() {
			welcome.SetText("Hello, " + entry.Text + " 😀\nWelcome to the Tingo/Tamago Fyne app!")
		}))

	w.SetContent(content)
	w.Canvas().Focus(entry)
}
