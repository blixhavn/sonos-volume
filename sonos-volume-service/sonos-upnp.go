package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

type Root struct {
	Device Device `xml:"device"`
}

type Device struct {
	RoomName string `xml:"roomName"`
}

func getRoomName(location string) (string, error) {
	resp, err := http.Get(location)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var root Root
	if err := xml.NewDecoder(resp.Body).Decode(&root); err != nil {
		return "", err
	}

	return root.Device.RoomName, nil
}

func main() {
	var preferredRoomName string
	flag.StringVar(&preferredRoomName, "room-name", "", "the preferred room name")
	flag.Parse()

	if preferredRoomName == "" {
		fmt.Println("Usage: sonos-volume-service --room-name=<preferred-room-name>")
		return
	}

	// Discover UPnP devices on the network
	devices, err := goupnp.DiscoverDevices("urn:schemas-upnp-org:device:ZonePlayer:1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Select the Sonos speakers based on the room
	var selectedDevice *goupnp.RootDevice
	for _, device := range devices {

		roomName, err := getRoomName(device.Location.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		if roomName == preferredRoomName {
			fmt.Println("Preferred room found:", roomName)
			selectedDevice = device.Root
			break
		}

	}

	if selectedDevice == nil {
		fmt.Println("No Sonos speakers found in the specified room")
		return
	}

	// Get the RenderingControl service
	renderingControlService, err := av1.NewRenderingControl1ClientsFromRootDevice(selectedDevice, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the current volume
	volume, err := renderingControlService[0].GetVolume(0, "Master")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Current Volume:", volume)
}
