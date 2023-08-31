package main

import (
	"flag"
	"fmt"
)

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

	volume, err := sonos.GetVolume()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Current Volume:", volume)

	err = sonos.IncreaseVolume()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Volume increased")

	err = sonos.Play()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Playback started")

	// ...
}
