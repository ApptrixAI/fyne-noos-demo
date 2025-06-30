package main

import (
	_ "embed"
	"image"
	"io"
	"log"
	"time"

	"github.com/usbarmory/go-boot/uefi"
	"github.com/usbarmory/go-boot/uefi/x64"

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
	handleError := func(err error) {
		log.Println("> ", err)
		time.Sleep(time.Second * 3)

		_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetWarm)
	}
	gop, err := x64.UEFI.Boot.GetGraphicsOutput()
	if err != nil {
		handleError(err)
		return
	}

	mode, _ := gop.GetMode()
	info, _ := mode.GetInfo()
	ww := uint64(info.HorizontalResolution)
	hh := uint64(info.VerticalResolution)

	if err = runApp(gop, ww, hh); err != nil {
		handleError(err)
		return
	}

	_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetShutdown)
}

func runApp(gop *uefi.GraphicsOutput, ww, hh uint64) error {
	var queue = make(chan noos.Event)
	buf := make([]byte, ww*hh*4)

	a := app.New()
	app.SetNoOSDriver(func(img image.Image) {
		err := refresh(img, gop, buf, ww, hh)
		if err != nil {
			fyne.LogError("refresh error", err)
		}
	}, queue)
	w := a.NewWindow("Fyne NOOS Demo")
	w.Resize(fyne.NewSize(float32(ww), float32(hh)))

	entry := widget.NewEntry()
	w.SetContent(makeUI(entry))
	w.Canvas().Focus(entry)

	w.Show()
	return runLoop(a, queue)
}

func makeUI(in fyne.CanvasObject) fyne.CanvasObject {
	welcome := widget.NewLabel("Hello NOOS Fyne!\n(press Ctrl+Q to quit)")
	welcome.Alignment = fyne.TextAlignCenter
	img := canvas.NewImageFromResource(fyne.NewStaticResource("fyne.png", imgData))
	img.SetMinSize(fyne.NewSquareSize(64))

	return container.NewVBox(welcome, container.NewCenter(img),
		widget.NewButton("Tap me", func() {
			welcome.SetText("Welcome 😀")
		}), in)
}

func refresh(img image.Image, gop *uefi.GraphicsOutput, buf []byte, ww, hh uint64) error {
	i := 0

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pix := img.At(x, y)

			r, g, b, _ := pix.RGBA()
			buf[i+0] = byte(b)
			buf[i+1] = byte(g)
			buf[i+2] = byte(r)
			buf[i+3] = 0

			i += 4
		}
	}

	return gop.Blt(buf, uefi.EfiBltBufferToVideo, 0, 0, 0, 0, ww, hh, 0)
}

func runLoop(a fyne.App, queue chan noos.Event) error {
	// we have to Run on a goroutine because the UEFI.Console is used on main in other code...
	wait := make(chan struct{})
	go func() {
		a.Driver().Run()
		wait <- struct{}{}
	}()

	defer func() {
		a.Driver().Quit()
		<-wait
	}()

	data := make([]byte, 4)
	for {
		n, err := x64.UEFI.Console.Read(data)
		if err == io.EOF {
			return err
		}
		if err != nil {
			fyne.LogError("failed to read", err)
			continue
		}

		for i := 0; i < n && i < len(data); i++ {
			switch data[i] {
			case 0:
				continue
			case 17: // ctrl+Q
				return nil
			}

			triggerKeys(data[i], queue)
		}
	}
}

func triggerKeys(key byte, queue chan noos.Event) {
	// visible characters
	if key >= ' ' && key <= '~' {
		queue <- &noos.CharacterEvent{Rune: rune(key)}
		return
	}

	// other keys
	name := fyne.KeyName(key)
	queue <- &noos.KeyEvent{Name: name, Direction: noos.KeyPressed}
	queue <- &noos.KeyEvent{Name: name, Direction: noos.KeyReleased}
}
