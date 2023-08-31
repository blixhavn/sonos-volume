package main

import (
	"flag"
	"fmt"

	"github.com/karalabe/hid"
)

// Vendor ID and Product ID of the Digispark ATTiny85
const VENDOR_ID = 0x16d0  // 5808
const PRODUCT_ID = 0x0753 // 1875

func main() {
	var preferredRoomName string
	flag.StringVar(&preferredRoomName, "room-name", "", "the preferred room name")
	flag.Parse()

	if preferredRoomName == "" {
		fmt.Println("Usage: sonos-volume-service --room-name=<preferred-room-name>")
		return
	}

	sonos, err := NewSonos(preferredRoomName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Enumerate all HID devices
	devices := hid.Enumerate(VENDOR_ID, PRODUCT_ID)
	if len(devices) == 0 {
		fmt.Println("No HID devices found")
		return
	}

	// Open the first HID device
	device, err := devices[0].Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer device.Close()

	// Read the USB data in a loop
	buf := make([]byte, 1)
	for {
		_, err := device.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		data := buf[0]
		if data == 'P' {
			fmt.Println("Play/Pause")
			err = sonos.PlayPause()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else if data == '+' {
			fmt.Println("Volume Up")
			err = sonos.IncreaseVolume()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else if data == '-' {
			fmt.Println("Volume Down")
			err = sonos.DecreaseVolume()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// ...
}
