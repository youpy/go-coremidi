package coremidi

import (
	"testing"
)

func TestNumberOfDevices(t *testing.T) {
	devices, err := AllDevices()
	numberOfDevices := len(devices)

	if err != nil {
		t.Fatalf("failed to get devices")
	}

	if numberOfDevices <= 0 {
		t.Fatalf("invalid number of devices")
	}
}

func TestEntities(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	entities, _ := device.Entities()

	if len(entities) <= 0 {
		t.Fatalf("invalid number of entities")
	}
}

func TestManufacturer(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	value := device.Manufacturer()

	if len(value) == 0 {
		t.Fatalf("invalid manufacturer")
	}
}

func TestName(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	value := device.Name()

	if len(value) == 0 {
		t.Fatalf("invalid name")
	}
}
