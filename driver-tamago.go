//go:build tamago

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
	"fyne.io/fyne/v2/driver/noos"
)

var buf []byte

func quit() {
	_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetShutdown)
}

func runApp(a fyne.App, queue chan noos.Event) {
	err := runLoop(a, queue)
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	log.Println("> ", err)
	time.Sleep(time.Second * 3)

	_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetWarm)
}

func refresh(img image.Image) {
	gop, _ := x64.UEFI.Boot.GetGraphicsOutput()
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

	_ = gop.Blt(buf, uefi.EfiBltBufferToVideo, 0, 0, 0, 0, uint64(width), uint64(height), 0)
}

func runLoop(a fyne.App, queue chan noos.Event) error {
	// we have to Run on a goroutine because the UEFI.Console is used on main in other code...
	wait := make(chan struct{})
	go func() {
		a.Run()
		wait <- struct{}{}
	}()

	defer func() {
		a.Quit()
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
			case 17, 23: // ctrl+Q or Esc
				return nil
			}

			triggerKey(uint16(data[i]), noos.KeyPressed, queue)
			triggerKey(uint16(data[i]), noos.KeyReleased, queue)
		}
	}
}

func screenSize() (uint64, uint64) {
	gop, err := x64.UEFI.Boot.GetGraphicsOutput()
	if err != nil {
		handleError(err)
		return 0, 0
	}

	mode, _ := gop.GetMode()
	info, _ := mode.GetInfo()
	ww, hh := uint64(info.HorizontalResolution),
		uint64(info.VerticalResolution)

	buf = make([]byte, ww*hh*4)
	return ww, hh
}
