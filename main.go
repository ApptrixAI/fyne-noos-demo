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

	setContent(w)
	w.Show()

	runApp(a, queue)
	quit()
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

func triggerKey(key uint16, dir noos.KeyDirection, queue chan noos.Event) {
	name := fyne.KeyName("")
	switch key {
	case 27: // esc
		name = fyne.KeyEscape
	case 13: // ret
		name = fyne.KeyReturn
	case 8: // backspace
		name = fyne.KeyBackspace
	case 9: // tab
		name = fyne.KeyTab
	case ' ':
		name = fyne.KeySpace
	default:
		if key > ' ' && key < '~' {
			name = fyne.KeyName(rune(key))
			if key >= 'a' && key <= 'z' {
				name = fyne.KeyName(rune(key) - 'a' + 'A')
			}
		}
	}
	queue <- &noos.KeyEvent{Name: name, Direction: dir}

	// visible characters
	if key >= ' ' && key <= '~' {
		if dir == noos.KeyReleased {
			return
		}

		queue <- &noos.CharacterEvent{Rune: rune(key)}
	}
}
