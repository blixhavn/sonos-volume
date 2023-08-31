package main

import (
	"encoding/xml"
	"errors"
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

type Sonos struct {
	avTransportService      *av1.AVTransport1
	renderingControlService *av1.RenderingControl1
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

func (s *Sonos) GetVolume() (uint16, error) {
	return s.renderingControlService.GetVolume(0, "Master")
}

func (s *Sonos) SetVolume(volume uint16) error {
	return s.renderingControlService.SetVolume(0, "Master", volume)
}

func (s *Sonos) IncreaseVolume() error {
	volume, err := s.GetVolume()
	if err != nil {
		return err
	}

	return s.SetVolume(volume + 1)
}

func (s *Sonos) DecreaseVolume() error {
	volume, err := s.GetVolume()
	if err != nil {
		return err
	}

	return s.SetVolume(volume - 1)
}

func (s *Sonos) Play() error {
	return s.avTransportService.Play(0, "1")
}

func (s *Sonos) Pause() error {
	return s.avTransportService.Pause(0)
}

func (s *Sonos) PlayPause() error {
	currentTransportState, _, _, err := s.avTransportService.GetTransportInfo(0)
	if err != nil {
		return err
	}

	if currentTransportState == "PLAYING" {
		return s.Pause()
	} else {
		return s.Play()
	}
}

func NewSonos(preferredRoomName string) (*Sonos, error) {

	// Discover UPnP devices on the network
	devices, err := goupnp.DiscoverDevices("urn:schemas-upnp-org:device:ZonePlayer:1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Select the Sonos speakers based on the room
	var selectedDevice *goupnp.RootDevice
	for _, device := range devices {

		roomName, err := getRoomName(device.Location.String())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if roomName == preferredRoomName {
			fmt.Println("Preferred room found:", roomName)
			selectedDevice = device.Root
			break
		}

	}

	if selectedDevice == nil {
		fmt.Println("No Sonos speakers found in the specified room")
		return nil, errors.New("No Sonos speakers found in the specified room")
	}

	avTransportServices, err := av1.NewAVTransport1ClientsFromRootDevice(selectedDevice, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Get the RenderingControl service
	renderingControlServices, err := av1.NewRenderingControl1ClientsFromRootDevice(selectedDevice, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Sonos{
		avTransportService:      avTransportServices[0],
		renderingControlService: renderingControlServices[0],
	}, nil
}
