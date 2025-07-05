package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/noos"
	"fyne.io/fyne/v2/widget"
)

//go:embed fyne.png
var imgData []byte

func main() {
	ww, hh := screenSize()
	var queue = make(chan noos.Event)

	a := app.New()
	app.SetNoOSDriver(refresh, queue)
	w := a.NewWindow("Fyne NOOS Demo")
	w.Resize(fyne.NewSize(float32(ww), float32(hh)))

	entry := widget.NewEntry()
	w.SetContent(makeUI(entry))
	w.Canvas().Focus(entry)

	w.Show()
	runApp(a, queue)
	quit()
}

func makeUI(in fyne.CanvasObject) fyne.CanvasObject {
	welcome := widget.NewLabel("Hello NoOS Fyne!\n(press Esc to quit)")
	welcome.Alignment = fyne.TextAlignCenter
	img := canvas.NewImageFromResource(fyne.NewStaticResource("fyne.png", imgData))
	img.SetMinSize(fyne.NewSquareSize(64))

	return container.NewVBox(welcome, container.NewCenter(img),
		widget.NewButton("Tap me", func() {
			welcome.SetText("Welcome 😀")
		}), in)
}

func triggerKey(key uint16, dir noos.KeyDirection, queue chan noos.Event) {
	// visible characters
	if key >= ' ' && key <= '~' {
		if dir == noos.KeyReleased {
			return
		}

		queue <- &noos.CharacterEvent{Rune: rune(key)}
		return
	}

	// other keys
	name := fyne.KeyName("")
	switch key {
	case 27: // esc
		name = fyne.KeyEscape
	case 13: // ret
		name = fyne.KeyReturn
	case 9: // backspace
		name = fyne.KeyBackspace
	case 8: // tab
		name = fyne.KeyTab
	default:
		name = fyne.KeyName([]byte{byte(key)})
	}
	queue <- &noos.KeyEvent{Name: name, Direction: dir}
}
