package coremidi

import (
	"errors"
	"testing"
)

func TestSources(t *testing.T) {
	device, err := findDeviceWithEntities()
	if err != nil {
		t.Fatal(err)
	}

	entities, _ := device.Entities()
	entity := entities[0]
	sources, _ := entity.Sources()

	if len(sources) <= 0 {
		t.Fatalf("invalid number of sources")
	}
}

func TestDestinations(t *testing.T) {
	device, err := findDeviceWithEntities()
	if err != nil {
		t.Fatal(err)
	}

	entities, _ := device.Entities()
	entity := entities[0]
	destinations, _ := entity.Destinations()

	if len(destinations) <= 0 {
		t.Fatalf("invalid number of destinations")
	}
}

func TestDevice(t *testing.T) {
	device, err := findDeviceWithEntities()
	if err != nil {
		t.Fatal(err)
	}

	entities, _ := device.Entities()
	entity := entities[0]
	result := entity.Device()

	if result.Name() != device.Name() {
		t.Fatalf("invalid name of device")
	}
}

func findDeviceWithEntities() (Device, error) {
	var device Device

	devices, _ := AllDevices()

	for _, dev := range devices {
		entities, _ := dev.Entities()

		if len(entities) > 0 {
			return dev, nil
		}
	}

	return device, errors.New("no device with entity was found")
}
