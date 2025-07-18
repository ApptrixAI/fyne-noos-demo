//go:build tinygo || noos

package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/noos"
)

const noKey = uint16(0xFFFF)

func quit() {
	// no-op
}

func runApp(a fyne.App, queue chan noos.Event) {
	go runEvents(a, queue)
	a.Run()
}

func runEvents(a fyne.App, queue chan noos.Event) {
	key := noKey
	for {
		time.Sleep(time.Millisecond * 10) // don't poll too fast

		newKey := nextKey()
		if newKey == key {
			continue
		}

		if newKey == noKey {
			triggerKey(key, noos.KeyReleased, queue)
			key = newKey
		} else {
			if newKey == 0x100 { // escape
				break
			}

			typed := mapKey(newKey)
			triggerKey(typed, noos.KeyPressed, queue)
			key = newKey
		}
	}

	fyne.Do(a.Quit)
}

func mapKey(key uint16) uint16 {
	// TODO handle shift...
	if key >= 'A' && key <= 'Z' {
		return key + 'a' - 'A'
	}

	switch key {
	case 0x100:
		return 27 // esc
	case 0x101:
		return 13 // ret
	case 0x102:
		return 8 // tab
	case 0x103:
		return 9 // backspace
	case 0x105:
		return 127 // delete
	}

	return key
}
